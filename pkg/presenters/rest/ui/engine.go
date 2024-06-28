package ui

import (
	"hash/fnv"
	"html/template"
)

func Unescape(s string) template.HTML {
	return template.HTML(s)
}
func NameToColor(s string) map[string]string {
	if s == "Unknown" {
		return map[string]string{
			"txt": "text-gray-600",
			"bg":  "bg-gray-200",
		}
	}
	colorPairs := []struct {
		text string
		bg   string
	}{
		{"text-red-600", "bg-red-200"},
		{"text-red-700", "bg-red-300"},
		{"text-red-800", "bg-red-400"},
		{"text-red-900", "bg-red-500"},
		{"text-pink-600", "bg-pink-200"},
		{"text-pink-700", "bg-pink-300"},
		{"text-pink-800", "bg-pink-400"},
		{"text-pink-900", "bg-pink-500"},
		{"text-rose-600", "bg-rose-200"},
		{"text-rose-700", "bg-rose-300"},
		{"text-rose-800", "bg-rose-400"},
		{"text-rose-900", "bg-rose-500"},
		{"text-yellow-600", "bg-yellow-200"},
		{"text-yellow-700", "bg-yellow-300"},
		{"text-yellow-800", "bg-yellow-400"},
		{"text-yellow-900", "bg-yellow-500"},
		{"text-fuchsia-600", "bg-fuchsia-200"},
		{"text-fuchsia-700", "bg-fuchsia-300"},
		{"text-fuchsia-800", "bg-fuchsia-400"},
		{"text-fuchsia-900", "bg-fuchsia-500"},
		{"text-lime-600", "bg-lime-200"},
		{"text-lime-700", "bg-lime-300"},
		{"text-lime-800", "bg-lime-400"},
		{"text-lime-900", "bg-lime-500"},
		{"text-green-600", "bg-green-200"},
		{"text-green-700", "bg-green-300"},
		{"text-green-800", "bg-green-400"},
		{"text-green-900", "bg-green-500"},
		{"text-teal-600", "bg-teal-200"},
		{"text-teal-700", "bg-teal-300"},
		{"text-teal-800", "bg-teal-400"},
		{"text-teal-900", "bg-teal-500"},
		{"text-sky-600", "bg-sky-200"},
		{"text-sky-700", "bg-sky-300"},
		{"text-sky-800", "bg-sky-400"},
		{"text-sky-900", "bg-sky-500"},
		{"text-azure-600", "bg-azure-200"},
		{"text-azure-700", "bg-azure-300"},
		{"text-azure-800", "bg-azure-400"},
		{"text-azure-900", "bg-azure-500"},
		{"text-indigo-600", "bg-indigo-200"},
		{"text-indigo-700", "bg-indigo-300"},
		{"text-indigo-800", "bg-indigo-400"},
		{"text-indigo-900", "bg-indigo-500"},
		{"text-blue-600", "bg-blue-200"},
		{"text-blue-700", "bg-blue-300"},
		{"text-blue-800", "bg-blue-400"},
		{"text-blue-900", "bg-blue-500"},
		{"text-emerald-600", "bg-emerald-200"},
		{"text-emerald-700", "bg-emerald-300"},
		{"text-emerald-800", "bg-emerald-400"},
		{"text-emerald-900", "bg-emerald-500"},
		{"text-cyan-600", "bg-cyan-200"},
		{"text-cyan-700", "bg-cyan-300"},
		{"text-cyan-800", "bg-cyan-400"},
		{"text-cyan-900", "bg-cyan-500"},
		{"text-violet-600", "bg-violet-200"},
		{"text-violet-700", "bg-violet-300"},
		{"text-violet-800", "bg-violet-400"},
		{"text-violet-900", "bg-violet-500"},
		{"text-lavender-600", "bg-lavender-200"},
		{"text-lavender-700", "bg-lavender-300"},
		{"text-lavender-800", "bg-lavender-400"},
		{"text-lavender-900", "bg-lavender-500"},
	}

	// Calculate a hash value for string s
	hash := fnv.New32a()
	hash.Write([]byte(s))
	hashValue := hash.Sum32()

	// Determine the index based on hash value
	index := int(hashValue) % len(colorPairs)

	// Retrieve the paired text and background colors
	colorPair := colorPairs[index]

	// Create a map to store text and background colors
	colors := map[string]string{
		"txt": colorPair.text,
		"bg":  colorPair.bg,
	}

	return colors
}
