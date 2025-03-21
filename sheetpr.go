// Copyright 2016 - 2025 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to and
// read from XLAM / XLSM / XLSX / XLTM / XLTX files. Supports reading and
// writing spreadsheet documents generated by Microsoft Excel™ 2007 and later.
// Supports complex components by high compatibility, and provided streaming
// API for generating or reading data from a worksheet with huge amounts of
// data. This library needs Go version 1.23 or later.

package excelize

import "reflect"

// SetPageMargins provides a function to set worksheet page margins.
func (f *File) SetPageMargins(sheet string, opts *PageLayoutMarginsOptions) error {
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	if opts == nil {
		return err
	}
	preparePageMargins := func(ws *xlsxWorksheet) {
		if ws.PageMargins == nil {
			ws.PageMargins = new(xlsxPageMargins)
		}
	}
	preparePrintOptions := func(ws *xlsxWorksheet) {
		if ws.PrintOptions == nil {
			ws.PrintOptions = new(xlsxPrintOptions)
		}
	}
	s := reflect.ValueOf(opts).Elem()
	for i := 0; i < 6; i++ {
		if !s.Field(i).IsNil() {
			preparePageMargins(ws)
			name := s.Type().Field(i).Name
			reflect.ValueOf(ws.PageMargins).Elem().FieldByName(name).Set(s.Field(i).Elem())
		}
	}
	if opts.Horizontally != nil {
		preparePrintOptions(ws)
		ws.PrintOptions.HorizontalCentered = *opts.Horizontally
	}
	if opts.Vertically != nil {
		preparePrintOptions(ws)
		ws.PrintOptions.VerticalCentered = *opts.Vertically
	}
	return err
}

// GetPageMargins provides a function to get worksheet page margins.
func (f *File) GetPageMargins(sheet string) (PageLayoutMarginsOptions, error) {
	opts := PageLayoutMarginsOptions{
		Bottom: float64Ptr(0.75),
		Footer: float64Ptr(0.3),
		Header: float64Ptr(0.3),
		Left:   float64Ptr(0.7),
		Right:  float64Ptr(0.7),
		Top:    float64Ptr(0.75),
	}
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return opts, err
	}
	if ws.PageMargins != nil {
		opts.Bottom = float64Ptr(ws.PageMargins.Bottom)
		opts.Footer = float64Ptr(ws.PageMargins.Footer)
		opts.Header = float64Ptr(ws.PageMargins.Header)
		opts.Left = float64Ptr(ws.PageMargins.Left)
		opts.Right = float64Ptr(ws.PageMargins.Right)
		opts.Top = float64Ptr(ws.PageMargins.Top)
	}
	if ws.PrintOptions != nil {
		opts.Horizontally = boolPtr(ws.PrintOptions.HorizontalCentered)
		opts.Vertically = boolPtr(ws.PrintOptions.VerticalCentered)
	}
	return opts, err
}

// prepareSheetPr create sheetPr element which not exist.
func (ws *xlsxWorksheet) prepareSheetPr() {
	if ws.SheetPr == nil {
		ws.SheetPr = new(xlsxSheetPr)
	}
}

// setSheetOutlineProps set worksheet outline properties by given options.
func (ws *xlsxWorksheet) setSheetOutlineProps(opts *SheetPropsOptions) {
	prepareOutlinePr := func(ws *xlsxWorksheet) {
		ws.prepareSheetPr()
		if ws.SheetPr.OutlinePr == nil {
			ws.SheetPr.OutlinePr = new(xlsxOutlinePr)
		}
	}
	if opts.OutlineSummaryBelow != nil {
		prepareOutlinePr(ws)
		ws.SheetPr.OutlinePr.SummaryBelow = opts.OutlineSummaryBelow
	}
	if opts.OutlineSummaryRight != nil {
		prepareOutlinePr(ws)
		ws.SheetPr.OutlinePr.SummaryRight = opts.OutlineSummaryRight
	}
}

// setSheetProps set worksheet format properties by given options.
func (ws *xlsxWorksheet) setSheetProps(opts *SheetPropsOptions) {
	preparePageSetUpPr := func(ws *xlsxWorksheet) {
		ws.prepareSheetPr()
		if ws.SheetPr.PageSetUpPr == nil {
			ws.SheetPr.PageSetUpPr = new(xlsxPageSetUpPr)
		}
	}
	prepareTabColor := func(ws *xlsxWorksheet) {
		ws.prepareSheetPr()
		if ws.SheetPr.TabColor == nil {
			ws.SheetPr.TabColor = new(xlsxColor)
		}
	}
	if opts.CodeName != nil {
		ws.prepareSheetPr()
		ws.SheetPr.CodeName = *opts.CodeName
	}
	if opts.EnableFormatConditionsCalculation != nil {
		ws.prepareSheetPr()
		ws.SheetPr.EnableFormatConditionsCalculation = opts.EnableFormatConditionsCalculation
	}
	if opts.Published != nil {
		ws.prepareSheetPr()
		ws.SheetPr.Published = opts.Published
	}
	if opts.AutoPageBreaks != nil {
		preparePageSetUpPr(ws)
		ws.SheetPr.PageSetUpPr.AutoPageBreaks = *opts.AutoPageBreaks
	}
	if opts.FitToPage != nil {
		preparePageSetUpPr(ws)
		ws.SheetPr.PageSetUpPr.FitToPage = *opts.FitToPage
	}
	ws.setSheetOutlineProps(opts)
	s := reflect.ValueOf(opts).Elem()
	for i := 5; i < 9; i++ {
		if !s.Field(i).IsNil() {
			prepareTabColor(ws)
			name := s.Type().Field(i).Name
			fld := reflect.ValueOf(ws.SheetPr.TabColor).Elem().FieldByName(name[8:])
			if s.Field(i).Kind() == reflect.Ptr && fld.Kind() == reflect.Ptr {
				fld.Set(s.Field(i))
				continue
			}
			fld.Set(s.Field(i).Elem())
		}
	}
}

// SetSheetProps provides a function to set worksheet properties.
func (f *File) SetSheetProps(sheet string, opts *SheetPropsOptions) error {
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return err
	}
	if opts == nil {
		return err
	}
	ws.setSheetProps(opts)
	if ws.SheetFormatPr == nil {
		ws.SheetFormatPr = &xlsxSheetFormatPr{DefaultRowHeight: defaultRowHeight}
	}
	s := reflect.ValueOf(opts).Elem()
	for i := 11; i < 18; i++ {
		if !s.Field(i).IsNil() {
			name := s.Type().Field(i).Name
			reflect.ValueOf(ws.SheetFormatPr).Elem().FieldByName(name).Set(s.Field(i).Elem())
		}
	}
	return err
}

// GetSheetProps provides a function to get worksheet properties.
func (f *File) GetSheetProps(sheet string) (SheetPropsOptions, error) {
	baseColWidth := uint8(8)
	opts := SheetPropsOptions{
		EnableFormatConditionsCalculation: boolPtr(true),
		Published:                         boolPtr(true),
		AutoPageBreaks:                    boolPtr(true),
		OutlineSummaryBelow:               boolPtr(true),
		BaseColWidth:                      &baseColWidth,
	}
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return opts, err
	}
	if ws.SheetPr != nil {
		opts.CodeName = stringPtr(ws.SheetPr.CodeName)
		if ws.SheetPr.EnableFormatConditionsCalculation != nil {
			opts.EnableFormatConditionsCalculation = ws.SheetPr.EnableFormatConditionsCalculation
		}
		if ws.SheetPr.Published != nil {
			opts.Published = ws.SheetPr.Published
		}
		if ws.SheetPr.PageSetUpPr != nil {
			opts.AutoPageBreaks = boolPtr(ws.SheetPr.PageSetUpPr.AutoPageBreaks)
			opts.FitToPage = boolPtr(ws.SheetPr.PageSetUpPr.FitToPage)
		}
		if ws.SheetPr.OutlinePr != nil {
			opts.OutlineSummaryBelow = ws.SheetPr.OutlinePr.SummaryBelow
			opts.OutlineSummaryRight = ws.SheetPr.OutlinePr.SummaryRight
		}
		if ws.SheetPr.TabColor != nil {
			opts.TabColorIndexed = intPtr(ws.SheetPr.TabColor.Indexed)
			opts.TabColorRGB = stringPtr(ws.SheetPr.TabColor.RGB)
			opts.TabColorTheme = ws.SheetPr.TabColor.Theme
			opts.TabColorTint = float64Ptr(ws.SheetPr.TabColor.Tint)
		}
	}
	if ws.SheetFormatPr != nil {
		opts.BaseColWidth = &ws.SheetFormatPr.BaseColWidth
		opts.DefaultColWidth = float64Ptr(ws.SheetFormatPr.DefaultColWidth)
		opts.DefaultRowHeight = float64Ptr(ws.SheetFormatPr.DefaultRowHeight)
		opts.CustomHeight = boolPtr(ws.SheetFormatPr.CustomHeight)
		opts.ZeroHeight = boolPtr(ws.SheetFormatPr.ZeroHeight)
		opts.ThickTop = boolPtr(ws.SheetFormatPr.ThickTop)
		opts.ThickBottom = boolPtr(ws.SheetFormatPr.ThickBottom)
	}
	return opts, err
}
