package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Unit represents a unit with its conversion factor to the base unit and its dimension.
type Unit struct {
	Factor    float64 // Factor to convert to the base unit
	Dimension string  // e.g., "mass" or "length"
	Name      string  // Full name of the unit
	// For temperature conversions, we need offset besides the factor
	Offset float64 // Used primarily for temperature conversions
}

// ConversionResult represents the result of a conversion operation
type ConversionResult struct {
	Success         bool    `json:"success"`
	Result          float64 `json:"result,omitempty"`
	FormattedResult string  `json:"formattedResult,omitempty"`
	Error           string  `json:"error,omitempty"`
	FromUnit        string  `json:"fromUnit,omitempty"`
	ToUnit          string  `json:"toUnit,omitempty"`
	InputValue      float64 `json:"inputValue,omitempty"`
}

// UnitConverter contains a mapping of unit symbols to their definitions.
type UnitConverter struct {
	units map[string]Unit
}

// NewUnitConverter initializes the converter with all unit dimensions.
func NewUnitConverter() *UnitConverter {
	return &UnitConverter{
		units: map[string]Unit{
			// Mass units (base = gram)
			"mg": {Factor: 0.001, Dimension: "mass", Name: "Milligram"},
			"g":  {Factor: 1, Dimension: "mass", Name: "Gram"},
			"kg": {Factor: 1000, Dimension: "mass", Name: "Kilogram"},
			"t":  {Factor: 1000000, Dimension: "mass", Name: "Tonne"},
			"oz": {Factor: 28.3495, Dimension: "mass", Name: "Ounce"},
			"lb": {Factor: 453.59237, Dimension: "mass", Name: "Pound"},

			// Length units (base = meter)
			"nm": {Factor: 0.000000001, Dimension: "length", Name: "Nanometer"},
			"µm": {Factor: 0.000001, Dimension: "length", Name: "Micrometer"},
			"mm": {Factor: 0.001, Dimension: "length", Name: "Millimeter"},
			"cm": {Factor: 0.01, Dimension: "length", Name: "Centimeter"},
			"m":  {Factor: 1, Dimension: "length", Name: "Meter"},
			"km": {Factor: 1000, Dimension: "length", Name: "Kilometer"},
			"in": {Factor: 0.0254, Dimension: "length", Name: "Inch"},
			"ft": {Factor: 0.3048, Dimension: "length", Name: "Foot"},
			"yd": {Factor: 0.9144, Dimension: "length", Name: "Yard"},
			"mi": {Factor: 1609.344, Dimension: "length", Name: "Mile"},

			// Temperature units (base = Kelvin)
			// For temperature, we need both factor and offset
			"C":  {Factor: 1, Offset: 273.15, Dimension: "temperature", Name: "Celsius"},
			"F":  {Factor: 5.0 / 9.0, Offset: 255.372, Dimension: "temperature", Name: "Fahrenheit"},
			"K":  {Factor: 1, Offset: 0, Dimension: "temperature", Name: "Kelvin"},
			"Ra": {Factor: 5.0 / 9.0, Offset: 0, Dimension: "temperature", Name: "Rankine"},

			// Time units (base = second)
			"ns":   {Factor: 1e-9, Dimension: "time", Name: "Nanosecond"},
			"µs":   {Factor: 1e-6, Dimension: "time", Name: "Microsecond"},
			"ms":   {Factor: 1e-3, Dimension: "time", Name: "Millisecond"},
			"s":    {Factor: 1, Dimension: "time", Name: "Second"},
			"min":  {Factor: 60, Dimension: "time", Name: "Minute"},
			"h":    {Factor: 3600, Dimension: "time", Name: "Hour"},
			"day":  {Factor: 86400, Dimension: "time", Name: "Day"},
			"week": {Factor: 604800, Dimension: "time", Name: "Week"},
			"year": {Factor: 31536000, Dimension: "time", Name: "Year (365 days)"},

			// Frequency units (base = hertz)
			"Hz":  {Factor: 1, Dimension: "frequency", Name: "Hertz"},
			"kHz": {Factor: 1000, Dimension: "frequency", Name: "Kilohertz"},
			"MHz": {Factor: 1e6, Dimension: "frequency", Name: "Megahertz"},
			"GHz": {Factor: 1e9, Dimension: "frequency", Name: "Gigahertz"},
			"THz": {Factor: 1e12, Dimension: "frequency", Name: "Terahertz"},

			// Speed units (base = meters per second)
			"m/s":  {Factor: 1, Dimension: "speed", Name: "Meters per second"},
			"km/h": {Factor: 0.277778, Dimension: "speed", Name: "Kilometers per hour"},
			"ft/s": {Factor: 0.3048, Dimension: "speed", Name: "Feet per second"},
			"mph":  {Factor: 0.44704, Dimension: "speed", Name: "Miles per hour"},
			"knot": {Factor: 0.514444, Dimension: "speed", Name: "Knot"},
			"mach": {Factor: 340.29, Dimension: "speed", Name: "Mach (at sea level)"},

			// Volume units (base = cubic meter)
			"m³":    {Factor: 1, Dimension: "volume", Name: "Cubic Meter"},
			"L":     {Factor: 0.001, Dimension: "volume", Name: "Liter"},
			"gal":   {Factor: 0.003785411784, Dimension: "volume", Name: "Gallon (US)"},
			"fl_oz": {Factor: 0.0000295735295625, Dimension: "volume", Name: "Fluid Ounce (US)"},

			// Area units (base = square meter)
			"m²":   {Factor: 1, Dimension: "area", Name: "Square Meter"},
			"acre": {Factor: 4046.8564224, Dimension: "area", Name: "Acre"},
			"ha":   {Factor: 10000, Dimension: "area", Name: "Hectare"},

			// Energy units (base = joule)
			"J":    {Factor: 1, Dimension: "energy", Name: "Joule"},
			"cal":  {Factor: 4.184, Dimension: "energy", Name: "Calorie"},
			"kcal": {Factor: 4184, Dimension: "energy", Name: "Kilocalorie"},

			// Power units (base = watt)
			"W":  {Factor: 1, Dimension: "power", Name: "Watt"},
			"HP": {Factor: 735.49875, Dimension: "power", Name: "Horsepower"},

			// Force units (base = newton)
			"N":   {Factor: 1, Dimension: "force", Name: "Newton"},
			"lbf": {Factor: 4.4482216153, Dimension: "force", Name: "Pound-force"},

			// Pressure units (base = pascal)
			"Pa":  {Factor: 1, Dimension: "pressure", Name: "Pascal"},
			"atm": {Factor: 101325, Dimension: "pressure", Name: "Atmosphere"},
			"bar": {Factor: 100000, Dimension: "pressure", Name: "Bar"},

			// Data Storage units (base = byte)
			"B":   {Factor: 1, Dimension: "data_storage", Name: "Byte"},
			"bit": {Factor: 0.125, Dimension: "data_storage", Name: "Bit"},
			"KB":  {Factor: 1024, Dimension: "data_storage", Name: "Kilobyte"},
			"MB":  {Factor: 1048576, Dimension: "data_storage", Name: "Megabyte"},
			"GB":  {Factor: 1073741824, Dimension: "data_storage", Name: "Gigabyte"},

			// Angle units (base = radian)
			"rad":    {Factor: 1, Dimension: "angle", Name: "Radian"},
			"deg":    {Factor: math.Pi / 180, Dimension: "angle", Name: "Degree"},
			"arcmin": {Factor: math.Pi / 10800, Dimension: "angle", Name: "Minute"},
			"arcsec": {Factor: math.Pi / 648000, Dimension: "angle", Name: "Second"},
		},
	}
}

// Convert performs the conversion from one unit to another.
func (uc *UnitConverter) Convert(value float64, from, to string) (float64, error) {
	unitFrom, ok := uc.units[from]
	if !ok {
		return 0, fmt.Errorf("invalid source unit: %s", from)
	}
	unitTo, ok := uc.units[to]
	if !ok {
		return 0, fmt.Errorf("invalid target unit: %s", to)
	}
	if unitFrom.Dimension != unitTo.Dimension {
		return 0, fmt.Errorf("cannot convert between different dimensions: %s (%s) and %s (%s)",
			from, unitFrom.Dimension, to, unitTo.Dimension)
	}

	var result float64

	// Special case for temperature, which needs offset handling
	if unitFrom.Dimension == "temperature" {
		// First convert to base unit (Kelvin)
		baseValue := value*unitFrom.Factor + unitFrom.Offset
		// Then convert from base unit to target unit
		result = (baseValue - unitTo.Offset) / unitTo.Factor
	} else {
		// Standard conversion via the base unit for other dimensions
		result = value * unitFrom.Factor / unitTo.Factor
	}

	// Round to a reasonable number of significant digits to avoid floating point issues
	precision := 12 // High precision to avoid truncation
	factor := math.Pow(10, float64(precision))
	return math.Round(result*factor) / factor, nil
}

// FormatResult formats the conversion result appropriately based on its magnitude
func (uc *UnitConverter) FormatResult(result float64, unit string) string {
	// Use scientific notation for very large or very small numbers
	absResult := math.Abs(result)
	if absResult < 0.001 || absResult > 1000000 {
		return fmt.Sprintf("%.6e %s", result, unit)
	}

	// For "normal" sized numbers, use appropriate decimal places
	var decimalPlaces int
	switch {
	case absResult >= 1000:
		decimalPlaces = 0
	case absResult >= 100:
		decimalPlaces = 1
	case absResult >= 10:
		decimalPlaces = 2
	case absResult >= 1:
		decimalPlaces = 3
	default:
		// For values less than 1, use more decimal places
		decimalPlaces = 4
	}

	formatString := fmt.Sprintf("%%.%df %%s", decimalPlaces)
	return fmt.Sprintf(formatString, result, unit)
}

// GetUnitsByDimension returns all units of a specific dimension
func (uc *UnitConverter) GetUnitsByDimension(dimension string) map[string]Unit {
	result := make(map[string]Unit)
	for symbol, unit := range uc.units {
		if unit.Dimension == dimension {
			result[symbol] = unit
		}
	}
	return result
}

// GetAllDimensions returns a slice of all available dimensions
func (uc *UnitConverter) GetAllDimensions() []string {
	dimensionMap := make(map[string]bool)
	for _, unit := range uc.units {
		dimensionMap[unit.Dimension] = true
	}

	dimensions := make([]string, 0, len(dimensionMap))
	for dim := range dimensionMap {
		dimensions = append(dimensions, dim)
	}
	return dimensions
}

// GetDimensionName returns a human-friendly name for a dimension
func (uc *UnitConverter) GetDimensionName(dimension string) string {
	switch dimension {
	case "volume":
		return "Volume"
	case "area":
		return "Area"
	case "energy":
		return "Energy"
	case "power":
		return "Power"
	case "force":
		return "Force"
	case "pressure":
		return "Pressure"
	case "data_storage":
		return "Data Storage"
	case "angle":
		return "Angle"
	case "mass":
		return "Mass"
	case "length":
		return "Length"
	case "temperature":
		return "Temperature"
	case "time":
		return "Time"
	case "frequency":
		return "Frequency"
	case "speed":
		return "Speed"
	default:
		return dimension
	}
}

// TemplateData represents the data passed to the HTML template
type TemplateData struct {
	Units          map[string]map[string]Unit // Map of dimensions to units
	Dimensions     []string
	DimensionNames map[string]string
	CurrentYear    int
}

// Handler for the homepage
func homeHandler(uc *UnitConverter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Organize units by dimension
		unitsByDimension := make(map[string]map[string]Unit)
		for _, dim := range uc.GetAllDimensions() {
			unitsByDimension[dim] = uc.GetUnitsByDimension(dim)
		}

		// Create dimension names map
		dimensionNames := make(map[string]string)
		for _, dim := range uc.GetAllDimensions() {
			dimensionNames[dim] = uc.GetDimensionName(dim)
		}

		// Create template data
		data := TemplateData{
			Units:          unitsByDimension,
			Dimensions:     uc.GetAllDimensions(),
			DimensionNames: dimensionNames,
			CurrentYear:    time.Now().Year(),
		}

		// Load template from index.html file
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error loading template: %v", err)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
	}
}

// Handler for the unit info endpoint
func unitInfoHandler(uc *UnitConverter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unitSymbol := r.URL.Query().Get("unit")
		if unitSymbol == "" {
			http.Error(w, "Unit symbol is required", http.StatusBadRequest)
			return
		}

		unit, ok := uc.units[unitSymbol]
		if !ok {
			http.Error(w, "Invalid unit symbol", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol":    unitSymbol,
			"name":      unit.Name,
			"dimension": unit.Dimension,
			"factor":    unit.Factor,
		})
	}
}

// Handler for getting units by dimension
func unitsByDimensionHandler(uc *UnitConverter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dimension := r.URL.Query().Get("dimension")
		if dimension == "" {
			http.Error(w, "Dimension is required", http.StatusBadRequest)
			return
		}

		units := uc.GetUnitsByDimension(dimension)
		if len(units) == 0 {
			http.Error(w, "Invalid dimension", http.StatusBadRequest)
			return
		}

		// Convert to a format suitable for the frontend
		type UnitInfo struct {
			Symbol string `json:"symbol"`
			Name   string `json:"name"`
		}

		unitInfos := make([]UnitInfo, 0, len(units))
		for symbol, unit := range units {
			unitInfos = append(unitInfos, UnitInfo{
				Symbol: symbol,
				Name:   unit.Name,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(unitInfos)
	}
}

// Handler for the conversion endpoint
func convertHandler(uc *UnitConverter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate headers
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			result := ConversionResult{
				Success: false,
				Error:   "Method not allowed. Please use POST.",
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(result)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			result := ConversionResult{
				Success: false,
				Error:   "Error parsing form data",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result)
			return
		}

		valueStr := r.FormValue("value")
		fromUnit := r.FormValue("from")
		toUnit := r.FormValue("to")

		// Validate input
		if valueStr == "" || fromUnit == "" || toUnit == "" {
			result := ConversionResult{
				Success: false,
				Error:   "All fields (value, from, to) are required",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result)
			return
		}

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			result := ConversionResult{
				Success: false,
				Error:   "Invalid value: must be a number",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(result)
			return
		}

		// Perform the conversion
		result, err := uc.Convert(value, fromUnit, toUnit)
		if err != nil {
			errorResult := ConversionResult{
				Success: false,
				Error:   err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResult)
			return
		}

		// Return the result as plain text (e.g., "10.00 kg")
		fmt.Fprintf(w, "%.3f %s", result, toUnit)
	}
}

func main() {
	uc := NewUnitConverter()

	// Define handlers
	http.HandleFunc("/", homeHandler(uc))
	http.HandleFunc("/convert", convertHandler(uc))
	http.HandleFunc("/unit-info", unitInfoHandler(uc))
	http.HandleFunc("/units-by-dimension", unitsByDimensionHandler(uc))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Add basic middleware for logging
	loggedRouter := logMiddleware(http.DefaultServeMux)

	// Start server
	port := ":8080"
	log.Printf("Server started on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, loggedRouter))
}

// Basic logging middleware
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log after request is processed
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}
