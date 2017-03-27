package main

import (
	"fmt"
	"github.com/nilsmagnus/grib/data"
	"os"
	"encoding/json"
	"bytes"
	"flag"
)


func main() {
	//category := flag.Int("category", 0, "Category. Default is temperature") // temperature
	//product := flag.Int("product", 6, "Product. Default is temperature") // temperature
	filename := flag.String("file", "", "Grib filename")

	flag.Parse()

	gribFile, err := os.Open(*filename)

	if err != nil {
		fmt.Printf("\nFile [%s] not found.\n", *filename)
	}
	defer gribFile.Close()

	messages, err := data.ReadAllMessages(gribFile)
	for _, message := range messages {
		fmt.Println(data.ReadDataType(int(message.Section1.Type)))
	}

}

func export(m *data.Message) {
	templateNumber := int(m.Section4.ProductDefinitionTemplateNumber)
	template := m.Section4.ProductDefinitionTemplate
	category := int(template.ParameterCategory)
	number := int(template.ParameterNumber)

	d := make(map[string]interface{})

	d["type"] = data.ReadDataType(int(m.Section1.Type));
	d["template"] = data.ReadProductDefinitionTemplateNumber(templateNumber);
	d["category"] = data.ReadProductDisciplineParameters(templateNumber, category);
	d["parameter"] = data.ReadProductDisciplineCategoryParameters(templateNumber, category, number);
	d["grid"] = data.ReadGridDefinitionTemplateNumber(int(m.Section3.TemplateNumber));
	d["surface1"] = data.ReadSurfaceTypesUnits(int(m.Section4.ProductDefinitionTemplate.FirstSurface.Type));
	d["surface1value"] = m.Section4.ProductDefinitionTemplate.FirstSurface.Value;
	d["surface1scale"] = m.Section4.ProductDefinitionTemplate.FirstSurface.Scale;
	d["surface2"] = data.ReadSurfaceTypesUnits(int(m.Section4.ProductDefinitionTemplate.SecondSurface.Type));
	d["surface2value"] = m.Section4.ProductDefinitionTemplate.SecondSurface.Value;
	d["data"] = m.Section7.Data;

	for k, v := range m.Section3.Definition.Export() {
		d[k] = v
	}

	// json print
	js, _ := json.Marshal(d)
	var out bytes.Buffer
	json.Indent(&out, js, "", "\t")
	out.WriteTo(os.Stdout)
	fmt.Println("")
}