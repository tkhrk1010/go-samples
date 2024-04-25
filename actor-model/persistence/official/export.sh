#!/bin/bash
# chmod +x export.sh

go_output_file="go_files.txt"
test_output_file="test_files.txt"

# 既存の出力ファイルを削除
rm -f "$go_output_file" "$test_output_file"

find . -name "*.go" -print0 | while IFS= read -r -d '' file; do
  if [[ $file == *"_test.go" ]]; then
    echo "=== $file ===" >> "$test_output_file"
    cat "$file" >> "$test_output_file"
    echo "" >> "$test_output_file"
  else
    echo "=== $file ===" >> "$go_output_file"
    cat "$file" >> "$go_output_file"
    echo "" >> "$go_output_file"
  fi
done

echo "Go files exported to $go_output_file"
echo "Test files exported to $test_output_file"