package terminal

import (
	"fmt"
	"runtime"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"temi/pkg"
)

/**

Implement terminal-based controller

xterm color reference https://jonasjacek.github.io/colors/
*/

const (
	// terminalWidth, width of terminal UI
	terminalWidth     = 120
	heapAllocBarCount = 6
)

// controller terminal-based controller
type controller struct {
	Grid *ui.Grid

	HeapObjectsSparkline     *widgets.Sparkline
	HeapObjectSparklineGroup *widgets.SparklineGroup
	HeapObjectsData          *pkg.StatRing

	SysText       *widgets.Paragraph
	GCCPUFraction *widgets.Gauge

	HeapAllocBarChart     *widgets.BarChart
	HeapAllocBarChartData *pkg.StatRing

	HeapPie *widgets.PieChart
}

func (p *controller) Resize() {
	p.resize()
	ui.Render(p.Grid)
}

func (p *controller) resize() {
	_, h := ui.TerminalDimensions()
	p.Grid.SetRect(0, 0, terminalWidth, h)
}

func (p *controller) Render(data *runtime.MemStats) {
	p.HeapObjectsData.Push(data.HeapObjects)
	p.HeapObjectsSparkline.Data = p.HeapObjectsData.NormalizedData()
	p.HeapObjectSparklineGroup.Title = fmt.Sprintf("HeapObjects, live heap object count: %d", data.HeapObjects)

	p.SysText.Text = fmt.Sprint(byteCountBinary(data.Sys))

	fNormalize := func() int {
		f := data.GCCPUFraction
		if f < 0.01 {
			for f < 1 {
				f = f * 10.0
			}
		}
		return int(f)
	}
	p.GCCPUFraction.Percent = fNormalize()
	p.GCCPUFraction.Label = fmt.Sprintf("%.2f%%", data.GCCPUFraction*100)

	p.HeapAllocBarChartData.Push(data.HeapAlloc)
	p.HeapAllocBarChart.Data = p.HeapAllocBarChartData.Data()
	p.HeapAllocBarChart.Labels = nil
	for _, v := range p.HeapAllocBarChart.Data {
		p.HeapAllocBarChart.Labels = append(p.HeapAllocBarChart.Labels, byteCountBinary(uint64(v)))
	}

	p.HeapPie.Data = []float64{float64(data.HeapIdle), float64(data.HeapInuse)}

	ui.Render(p.Grid)
}

func (p *controller) initUI() {
	p.resize()

	p.HeapObjectsSparkline.LineColor = ui.Color(89) // xterm color DeepPink4
	p.HeapObjectSparklineGroup = widgets.NewSparklineGroup(p.HeapObjectsSparkline)

	p.SysText.Title = "Sys, the total bytes of memory obtained from the OS"
	p.SysText.PaddingLeft = 25
	p.SysText.PaddingTop = 1

	p.HeapAllocBarChart.BarGap = 2
	p.HeapAllocBarChart.BarWidth = 8
	p.HeapAllocBarChart.Title = "HeapAlloc, bytes of allocated heap objects"
	p.HeapAllocBarChart.NumFormatter = func(f float64) string { return "" }

	p.GCCPUFraction.Title = "GCCPUFraction 0%~100%"
	p.GCCPUFraction.BarColor = ui.Color(50) // xterm color Cyan2

	p.HeapPie.Title = "HeapInuse vs HeapIdle"
	p.HeapPie.LabelFormatter = func(idx int, _ float64) string { return []string{"Idle", "Inuse"}[idx] }

	p.Grid.Set(
		ui.NewRow(.2, p.HeapObjectSparklineGroup),
		ui.NewRow(.8,
			ui.NewCol(.5,
				ui.NewRow(.2, p.SysText),
				ui.NewRow(.2, p.GCCPUFraction),
				ui.NewRow(.6, p.HeapAllocBarChart),
			),
			ui.NewCol(.5, p.HeapPie),
		),
	)

}

func newController() *controller {

	ctl := &controller{
		Grid: ui.NewGrid(),

		HeapObjectsSparkline: widgets.NewSparkline(),
		HeapObjectsData:      pkg.NewChartRing(terminalWidth),

		SysText:       widgets.NewParagraph(),
		GCCPUFraction: widgets.NewGauge(),

		HeapAllocBarChart:     widgets.NewBarChart(),
		HeapAllocBarChartData: pkg.NewChartRing(heapAllocBarCount),

		HeapPie: widgets.NewPieChart(),
	}

	ctl.initUI()

	return ctl
}
