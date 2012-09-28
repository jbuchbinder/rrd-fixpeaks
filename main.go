// RRD-FIXPEAKS
// https://github.com/jbuchbinder/rrd-fixpeaks

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"math"
	"os/exec"
	"strconv"
)

var (
	dryRun     = flag.Bool("dryrun", false, "Dry run flag (don't write)")
	threshold  = flag.Float64("threshold", 0, "Threshold percentage above avg above which values should be clipped")
	multiplier = flag.Float64("multiplier", 2, "Factor which max must outstrip average")
	minDiff    = flag.Float64("mindiff", 0, "Minimum difference above average")
	absAbove   = flag.Float64("absabove", -1, "If not -1, every value above this will be removed")
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: rrd-fixpeaks -threshold=80 -mindiff=10 -multiplier=1 -absabove=-1 RRDFILE.rrd")
		return
	}
	rrdfiles := flag.Args()

	for i := 0; i < len(rrdfiles); i++ {
		r := Rrd{}
		x := dumpXml(rrdfiles[i])
		err := xml.Unmarshal(x, &r)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", rrdfiles[i], err)
			continue
		}
		rrdInfo(rrdfiles[i], r)

		width := len(r.Rra[0].Database.Data[0].Value)
		allMax := float64(0)
		rMax := make([]float64, width)
		rAvg := make([]float64, width)
		rTotal := make([]float64, width)

		for j := 0; j < len(r.Rra); j++ {
			rraCount := len(r.Rra[j].Database.Data)
			for k := 0; k < rraCount; k++ {
				for l := 0; l < width; l++ {
					if r.Rra[j].Database.Data[k].Value[l] == "" || r.Rra[j].Database.Data[k].Value[l] == "NaN" {
						continue
					}
					v, err := strconv.ParseFloat(r.Rra[j].Database.Data[k].Value[l], 64)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					rTotal[l] += v
					if v > allMax {
						allMax = v
					}
					if v > rMax[l] {
						rMax[l] = v
					}
				}
			}
		}

		for j := 0; j < len(r.Rra); j++ {
			modCount := 0
			rraCount := len(r.Rra[j].Database.Data)

			rT := make([]float64, width)

			for m := 0; m < width; m++ {
				fmt.Printf("RRA #%d element %d max = %f allMax = %f\n", j, m, rMax[m], allMax)
				rAvg[m] = rTotal[m] / float64(rraCount)
				fmt.Printf("RRA #%d element %d average = %f\n", j, m, rAvg[m])
				mx := allMax /* rMax[m] */
				f := rAvg[m]
				fmt.Printf("starting = %f, floor = %f, tperc = %f\n", rMax[m]-f, f, (*threshold / float64(100)))
				rT[m] = ((mx - f) * (*threshold / float64(100))) + f
				fmt.Printf("RRA #%d element %d thold/floor = %f\n", j, m, rT[m])
			}

			// Second loop, correct values
			for k := 0; k < rraCount; k++ {
				modify := false
				for l := 0; l < width; l++ {
					if r.Rra[j].Database.Data[k].Value[l] == "" || r.Rra[j].Database.Data[k].Value[l] == "NaN" {
						continue
					}
					v, err := strconv.ParseFloat(r.Rra[j].Database.Data[k].Value[l], 64)
					if err != nil {
						continue
					}
					// Stop div by zero
					if v > 0 {
						if *absAbove > -1 && v > *absAbove {
							modify = true
						}
						if v > rT[l] {
							if *minDiff == 0 || math.Abs(v-rAvg[l]) > *minDiff {
								if *multiplier == 0 || *multiplier <= 1 || allMax/v >= *multiplier {
									modify = true
								}
							}
						}
					}
				}
				if modify {
					for l := 0; l < width; l++ {
						r.Rra[j].Database.Data[k].Value[l] = "NaN"
					}
					modCount++
				}
			}
			fmt.Printf("modcount = %d\n", modCount)
		}

		if !*dryRun {
			restoreXml(rrdfiles[i], r)
		}
	}

}

func rrdInfo(file string, rrd Rrd) {
	fmt.Printf("%s has %d RRAs\n", file, len(rrd.Rra))
	for i := 0; i < len(rrd.Rra); i++ {
		endTs := rrd.LastUpdate
		entries := len(rrd.Rra[i].Database.Data)
		incr := rrd.Step * rrd.Rra[i].PdpPerRow
		beginTs := endTs - ((int64(entries) - 1) * int64(incr))
		fmt.Printf("\t[%d] has %d entries\n", i, entries)
		fmt.Printf("\t\tRepresents %d sec increments (%d - %d)\n", incr, beginTs, endTs)
	}
}

func dumpXml(file string) []byte {
	out, err := exec.Command("rrdtool", "dump", file).Output()
	if err != nil {
		panic(err)
	}
	return out
}

func restoreXml(file string, rrd Rrd) {
	cmd := exec.Command("rrdtool", "restore", "-f", "-", file)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	bin, err := xml.Marshal(rrd)
	// DEBUG:
	// fmt.Println(string(bin))
	_, err = stdin.Write([]byte(bin))
	if err != nil {
		panic(err)
	}
	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
