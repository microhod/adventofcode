path = "input.txt"

{:ok, data} = File.read(path)

{first, second} =
  Enum.unzip(
    data
    |> String.split("\n", trim: true)
    |> Enum.map(fn line ->
      {first, line} = Integer.parse(line)
      line = String.trim(line)
      {second, ""} = Integer.parse(line)
      {first, second}
    end)
  )

# part 1
diff =
  Enum.zip([Enum.sort(first), Enum.sort(second)])
  |> Enum.map(fn {first, second} -> abs(first - second) end)
  |> Enum.sum()

IO.puts("Part 1:\ntotal diff: #{diff}")

# part 2
freq = Enum.frequencies(second)
score = Enum.reduce(first, fn n, acc -> acc + n * Map.get(freq, n, 0) end)

IO.puts("\nPart 2:\nsimilarity score: #{score}")
