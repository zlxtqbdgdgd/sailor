// Copyright 2018 ROOBO. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chart

import (
	"fmt"
	"os"

	"roobo.com/sailor/math"
)

var (
	svgHeader = "<?xml version='1.0' standalone='no'?>\n" +
		"<!DOCTYPE svg PUBLIC '-//W3C//DTD SVG 1.1//EN'\n" +
		"  'http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd'>\n" +
		"<svg width='%d' height='%d' viewBox='0 0 %[1]d %[2]d'\n" +
		"     xmlns='http://www.w3.org/2000/svg' version='1.1'>\n" +
		"  <desc>time cost(ms)</desc>\n" +
		"    <defs>\n" +
		"      <pattern id='smallGrid' width='10' height='10' patternUnits='userSpaceOnUse'>\n" +
		"        <path d='M 10 0 L 0 0 0 10' fill='none' stroke='gray' stroke-width='0.3'/>\n" +
		"      </pattern>\n" +
		"      <pattern id='grid' width='100' height='100' patternUnits='userSpaceOnUse'>\n" +
		"        <rect width='100' height='100' fill='url(#smallGrid)'/>\n" +
		"        <path d='M 100 0 L 0 0 0 100' fill='none' stroke='gray' stroke-width='1'/>\n" +
		"      </pattern>\n" +
		"    </defs>\n" +
		"    <rect width='100%%' height='100%%' fill='url(#grid)' />\n" +
		"    <rect x='0' y='0' width='%[1]d' height='%[2]d'\n" +
		"          fill='none' stroke='black' stroke-width='2' />\n"

	svgLine = "  <polyline fill='none' stroke='%[2]s' stroke-width='1'\n" +
		"          points='%[1]s' />\n"
	svgText = "  <text x='%d' y='%.2f'>%s</text>\n"
)

func WriteSvg(data []float64, filename, color string) error {
	svg := Polyline(data, color)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(svg); err != nil {
		fmt.Println("write string failed")
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// color: blue, red, brown,darkgreen,darkviolet,lime,navy,pink,seagreen,gold,
// tomato,yellowgreen, etc...
func Polyline(yaxis []float64, color string) string {
	max, idMax := math.MaxF64(yaxis...)
	min, idMin := math.MinF64(yaxis...)
	ave := math.AveF64(yaxis...)
	w := len(yaxis)
	h := int(max) * 10 / 8
	wfactor := 800.0 / float64(w)
	hfactor := 600.0 / float64(h)
	width, height := int(float64(w)*wfactor), int(float64(h)*hfactor)
	svg := fmt.Sprintf(svgHeader, width, height)
	data := ""
	for i := 0; i < len(yaxis); i++ {
		data += fmt.Sprintf(" %d,%.2f", int(float64(i)*wfactor),
			float64(height)-hfactor*yaxis[i])
		if i%5 == 0 && i != 0 {
			data += "\n                  "
		}
	}
	svg += fmt.Sprintf(svgLine, data, color)
	//
	idMax = int(float64(idMax) * wfactor)
	idMin = int(float64(idMin) * wfactor)
	svg += fmt.Sprintf(svgText, idMax, float64(height)-max*hfactor, fmt.Sprintf("max=%.2fms", max))
	svg += fmt.Sprintf(svgText, idMin, float64(height)-min*hfactor, fmt.Sprintf("min=%.2fms", min))
	svg += fmt.Sprintf(svgLine, fmt.Sprintf("0,%.2f %d,%.2[1]f", float64(height)-ave*hfactor, width), "blue")
	svg += fmt.Sprintf(svgText, 0, float64(height)-ave*hfactor, fmt.Sprintf("ave=%.2fms", ave))
	//
	svg += "</svg>"
	return svg
}
