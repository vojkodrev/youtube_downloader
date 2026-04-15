package main

type LowerQualityVideoRecoder interface {
	Recode(inputPath, outputPath string) error
}
