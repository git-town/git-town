. as $root
| (
  [ to_entries[]
    | [
        "\t",(.key|@json),": &GherkinDialect{\n",
        "\t\t", (.key|@json),", ", (.value.name|@json),", ", (.value.native|@json), ", map[string][]string{\n"
      ] + (
          [ .value
            | {"feature","rule","background","scenario","scenarioOutline","examples","given","when","then","and","but"}
            | to_entries[]
            | "\t\t\t"+(.key), ": []string{\n",
                ([ .value[] | "\t\t\t\t", @json, ",\n"  ]|add),
              "\t\t\t},\n"
          ]
      ) + ["\t\t},\n","\t},\n"]
    | add
  ]
  | add
  )
| "package gherkin\n\n"
+ "// Builtin dialects for " + ([ $root | to_entries[] | .key+" ("+.value.name+")" ] | join(", ")) + "\n"
+ "func GherkinDialectsBuildin() GherkinDialectProvider {\n"
+ "\treturn buildinDialects\n"
+ "}\n\n"
+ "const (\n"
+ (
  ["feature","rule","background","scenario","scenarioOutline","examples","given","when","then","and","but"]
  | [ .[] | "\t" + . + " = " + (.|@json) + "\n" ]
  | add )
+ ")\n\n"
+ "var buildinDialects = gherkinDialectMap{\n"
+ .
+ "}\n"