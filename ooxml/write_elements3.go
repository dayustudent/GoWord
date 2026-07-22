package ooxml

import (
	"fmt"
	"strings"

	"github.com/VantageDataChat/GoWord/common"
	"github.com/VantageDataChat/GoWord/document"
)

func (w *Writer) writeImage(b *strings.Builder, img *document.Image) {
	b.WriteString(`<w:p>`)
	switch img.Style.WrappingStyle {
	case "behind":
		w.writeBehindImage(b, img)
	case "inFrontOf":
		w.writeInFrontOfImage(b, img)
	case "square":
		w.writeSquareImage(b, img)
	case "tight":
		w.writeTightImage(b, img)
	default:
		w.writeInlineImage(b, img)
	}
	b.WriteString(`</w:p>`)
}

func (w *Writer) writeInlineImage(b *strings.Builder, img *document.Image) {
	width := img.Style.Width
	if width == 0 {
		width = 100
	}
	height := img.Style.Height
	if height == 0 {
		height = 100
	}

	cx := common.PointToEmu(width)
	cy := common.PointToEmu(height)

	b.WriteString(`<w:r>`)
	b.WriteString(`<w:drawing>`)
	b.WriteString(fmt.Sprintf(`<wp:inline distT="0" distB="0" distL="0" distR="0">`))
	b.WriteString(fmt.Sprintf(`<wp:extent cx="%d" cy="%d"/>`, cx, cy))
	b.WriteString(`<wp:docPr id="` + fmt.Sprintf("%d", img.ID) + `" name="` + escapeXML(img.Name) + `"/>`)
	b.WriteString(`<a:graphic>`)
	b.WriteString(fmt.Sprintf(`<a:graphicData uri="%s">`, nsPIC))
	b.WriteString(fmt.Sprintf(`<pic:pic xmlns:pic="%s">`, nsPIC))
	b.WriteString(fmt.Sprintf(`<pic:nvPicPr><pic:cNvPr id="%d" name="%s"/><pic:cNvPicPr/></pic:nvPicPr>`, img.ID, escapeXML(img.Name)))
	b.WriteString(fmt.Sprintf(`<pic:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>`, img.RelID))
	b.WriteString(fmt.Sprintf(`<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr>`, cx, cy))
	b.WriteString(`</pic:pic>`)
	b.WriteString(`</a:graphicData>`)
	b.WriteString(`</a:graphic>`)
	b.WriteString(`</wp:inline>`)
	b.WriteString(`</w:drawing>`)
	b.WriteString(`</w:r>`)
}

func (w *Writer) writeBehindImage(b *strings.Builder, img *document.Image) {
	w.writeFloatingImage(b, img, "1")
}

func (w *Writer) writeInFrontOfImage(b *strings.Builder, img *document.Image) {
	w.writeFloatingImage(b, img, "0")
}

func (w *Writer) writeSquareImage(b *strings.Builder, img *document.Image) {
	w.writeFloatingImage(b, img, "0", "square")
}

func (w *Writer) writeTightImage(b *strings.Builder, img *document.Image) {
	w.writeFloatingImage(b, img, "0", "tight")
}

func (w *Writer) writeFloatingImage(b *strings.Builder, img *document.Image, behindDoc string, wrapTypes ...string) {
	width := img.Style.Width
	if width == 0 {
		width = 100
	}
	height := img.Style.Height
	if height == 0 {
		height = 100
	}
	cx := common.PointToEmu(width)
	cy := common.PointToEmu(height)

	marginT := common.PointToEmu(img.Style.MarginTop)
	marginB := common.PointToEmu(img.Style.MarginBottom)
	marginL := common.PointToEmu(img.Style.MarginLeft)
	marginR := common.PointToEmu(img.Style.MarginRight)

	pageWidthEmu := int64(11907 * 635)
	pageHeightEmu := int64(16839 * 635)

	align := img.Style.Alignment
	if align == "" {
		align = "center"
	}

	posH := (pageWidthEmu - cx) / 2
	posV := (pageHeightEmu - cy) / 2
	if posH < 0 {
		posH = 0
	}
	if posV < 0 {
		posV = 0
	}

	if align == "left" {
		posH = 0
	} else if align == "right" {
		posH = pageWidthEmu - cx
		if posH < 0 {
			posH = 0
		}
	}
	if align == "top" {
		posV = 0
	} else if align == "bottom" {
		posV = pageHeightEmu - cy
		if posV < 0 {
			posV = 0
		}
	}

	wrapType := `<wp:wrapNone/>`
	if len(wrapTypes) > 0 {
		switch wrapTypes[0] {
		case "square":
			wrapType = `<wp:wrapSquare wrapText="both"/>`
		case "tight":
			wrapType = `<wp:wrapTight wrapText="both"/>`
		}
	}

	b.WriteString(`<w:r>`)
	b.WriteString(`<w:drawing>`)
	b.WriteString(fmt.Sprintf(`<wp:anchor distT="%d" distB="%d" distL="%d" distR="%d" simplePos="0" relativeHeight="251658240" behindDoc="%s" locked="0" layoutInCell="1" allowOverlap="1">`, marginT, marginB, marginL, marginR, behindDoc))
	b.WriteString(`<wp:simplePos x="0" y="0"/>`)
	b.WriteString(`<wp:positionH relativeFrom="page">`)
	b.WriteString(fmt.Sprintf(`<wp:posOffset>%d</wp:posOffset>`, posH))
	b.WriteString(`</wp:positionH>`)
	b.WriteString(`<wp:positionV relativeFrom="page">`)
	b.WriteString(fmt.Sprintf(`<wp:posOffset>%d</wp:posOffset>`, posV))
	b.WriteString(`</wp:positionV>`)
	b.WriteString(fmt.Sprintf(`<wp:extent cx="%d" cy="%d"/>`, cx, cy))
	b.WriteString(`<wp:effectExtent l="0" t="0" r="0" b="0"/>`)
	b.WriteString(wrapType)
	b.WriteString(fmt.Sprintf(`<wp:docPr id="%d" name="%s"/>`, img.ID, escapeXML(img.Name)))
	b.WriteString(`<wp:cNvGraphicFramePr>`)
	b.WriteString(`<a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/>`)
	b.WriteString(`</wp:cNvGraphicFramePr>`)
	b.WriteString(`<a:graphic>`)
	b.WriteString(fmt.Sprintf(`<a:graphicData uri="%s">`, nsPIC))
	b.WriteString(fmt.Sprintf(`<pic:pic xmlns:pic="%s">`, nsPIC))
	b.WriteString(fmt.Sprintf(`<pic:nvPicPr><pic:cNvPr id="%d" name="%s"/><pic:cNvPicPr/></pic:nvPicPr>`, img.ID, escapeXML(img.Name)))
	b.WriteString(fmt.Sprintf(`<pic:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>`, img.RelID))
	b.WriteString(fmt.Sprintf(`<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr>`, cx, cy))
	b.WriteString(`</pic:pic>`)
	b.WriteString(`</a:graphicData>`)
	b.WriteString(`</a:graphic>`)
	b.WriteString(`</wp:anchor>`)
	b.WriteString(`</w:drawing>`)
	b.WriteString(`</w:r>`)
}

func (w *Writer) writeTable(b *strings.Builder, t *document.Table) {
	b.WriteString(`<w:tbl>`)
	b.WriteString(`<w:tblPr>`)
	if t.StyleName != "" {
		b.WriteString(fmt.Sprintf(`<w:tblStyle w:val="%s"/>`, escapeXML(t.StyleName)))
	}
	ts := &t.Style
	if ts.Width > 0 {
		wType := ts.WidthType
		if wType == "" {
			wType = "dxa"
		}
		b.WriteString(fmt.Sprintf(`<w:tblW w:w="%d" w:type="%s"/>`, ts.Width, wType))
	} else {
		b.WriteString(`<w:tblW w:w="0" w:type="auto"/>`)
	}
	if ts.Alignment != "" {
		b.WriteString(fmt.Sprintf(`<w:jc w:val="%s"/>`, ts.Alignment))
	}
	w.writeTableBorders(b, ts)
	w.writeTableCellMargins(b, ts)
	if ts.Layout != "" {
		b.WriteString(fmt.Sprintf(`<w:tblLayout w:type="%s"/>`, ts.Layout))
	}
	if ts.CellSpacing > 0 {
		b.WriteString(fmt.Sprintf(`<w:tblCellSpacing w:w="%d" w:type="dxa"/>`, ts.CellSpacing))
	}
	if ts.Indent > 0 {
		b.WriteString(fmt.Sprintf(`<w:tblInd w:w="%d" w:type="dxa"/>`, ts.Indent))
	}
	b.WriteString(`</w:tblPr>`)
	if len(t.Grid) > 0 {
		b.WriteString(`<w:tblGrid>`)
		for _, colW := range t.Grid {
			b.WriteString(fmt.Sprintf(`<w:gridCol w:w="%d"/>`, colW))
		}
		b.WriteString(`</w:tblGrid>`)
	}
	for _, row := range t.Rows {
		w.writeTableRow(b, row)
	}
	b.WriteString(`</w:tbl>`)
}