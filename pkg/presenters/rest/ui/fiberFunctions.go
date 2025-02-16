package ui

import (
	"hash/fnv"
	"html/template"
)

func Unescape(s string) template.HTML {
	return template.HTML(s)
}

func UserInMap(expenseUsers map[string]struct {
	DisplayName      string
	TelegramUsername string
}, userID string) bool {
	_, exists := expenseUsers[userID]
	return exists
}

func NameToColor(s string) map[string]string {
	if s == "Unknown" {
		return map[string]string{
			"txt":    "text-gray-500",
			"bg":     "bg-gray-200",
			"border": "border-gray-200",
		}
	}
	colorPairs := []struct {
		text   string
		bg     string
		border string
	}{
		{"text-red-500", "bg-red-200", "border-red-300"},
		{"text-red-600", "bg-red-300", "border-red-400"},
		{"text-red-700", "bg-red-400", "border-red-500"},
		{"text-pink-500", "bg-pink-200", "border-pink-300"},
		{"text-pink-600", "bg-pink-300", "border-pink-400"},
		{"text-pink-700", "bg-pink-400", "border-pink-500"},
		{"text-pink-800", "bg-pink-500", "border-pink-600"},
		{"text-rose-500", "bg-rose-200", "border-rose-300"},
		{"text-rose-600", "bg-rose-300", "border-rose-400"},
		{"text-rose-700", "bg-rose-400", "border-rose-500"},
		{"text-rose-800", "bg-rose-500", "border-rose-600"},
		{"text-yellow-500", "bg-yellow-200", "border-yellow-300"},
		{"text-yellow-600", "bg-yellow-300", "border-yellow-400"},
		{"text-yellow-700", "bg-yellow-400", "border-yellow-500"},
		{"text-yellow-800", "bg-yellow-500", "border-yellow-600"},
		{"text-fuchsia-500", "bg-fuchsia-200", "border-fuchsia-300"},
		{"text-fuchsia-600", "bg-fuchsia-300", "border-fuchsia-400"},
		{"text-fuchsia-700", "bg-fuchsia-400", "border-fuchsia-500"},
		{"text-fuchsia-800", "bg-fuchsia-500", "border-fuchsia-600"},
		{"text-lime-500", "bg-lime-200", "border-lime-300"},
		{"text-lime-600", "bg-lime-300", "border-lime-400"},
		{"text-lime-700", "bg-lime-400", "border-lime-500"},
		{"text-lime-800", "bg-lime-500", "border-lime-600"},
		{"text-green-500", "bg-green-200", "border-green-300"},
		{"text-green-600", "bg-green-300", "border-green-400"},
		{"text-green-700", "bg-green-400", "border-green-500"},
		{"text-green-800", "bg-green-500", "border-green-600"},
		{"text-teal-500", "bg-teal-200", "border-teal-300"},
		{"text-teal-600", "bg-teal-300", "border-teal-400"},
		{"text-teal-700", "bg-teal-400", "border-teal-500"},
		{"text-teal-800", "bg-teal-500", "border-teal-600"},
		{"text-sky-500", "bg-sky-200", "border-sky-300"},
		{"text-sky-600", "bg-sky-300", "border-sky-400"},
		{"text-sky-700", "bg-sky-400", "border-sky-500"},
		{"text-sky-800", "bg-sky-500", "border-sky-600"},
		{"text-indigo-500", "bg-indigo-200", "border-indigo-300"},
		{"text-indigo-600", "bg-indigo-300", "border-indigo-400"},
		{"text-indigo-700", "bg-indigo-400", "border-indigo-500"},
		{"text-indigo-800", "bg-indigo-500", "border-indigo-600"},
		{"text-blue-500", "bg-blue-200", "border-blue-300"},
		{"text-blue-600", "bg-blue-300", "border-blue-400"},
		{"text-blue-700", "bg-blue-400", "border-blue-500"},
		{"text-blue-800", "bg-blue-500", "border-blue-600"},
		{"text-emerald-500", "bg-emerald-200", "border-emerald-300"},
		{"text-emerald-600", "bg-emerald-300", "border-emerald-400"},
		{"text-emerald-700", "bg-emerald-400", "border-emerald-500"},
		{"text-emerald-800", "bg-emerald-500", "border-emerald-600"},
		{"text-cyan-500", "bg-cyan-200", "border-cyan-300"},
		{"text-cyan-600", "bg-cyan-300", "border-cyan-400"},
		{"text-cyan-700", "bg-cyan-400", "border-cyan-500"},
		{"text-cyan-800", "bg-cyan-500", "border-cyan-600"},
		{"text-violet-600", "bg-violet-200", "border-violet-300"},
		{"text-violet-700", "bg-violet-300", "border-violet-400"},
		{"text-violet-800", "bg-violet-400", "border-violet-500"},
		{"text-violet-900", "bg-violet-500", "border-violet-600"},
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
		"txt":    colorPair.text,
		"bg":     colorPair.bg,
		"border": colorPair.border,
	}
	return colors
}
