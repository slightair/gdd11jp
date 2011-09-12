gets
gets

5.times do |i|
  Dir.glob("problems#{i+1}/*.txt") do |f|
    File.delete(f)
  end
end

count = 1
while line = gets
  line.chop!
  
  i = count / 1000 + 1
  i = 5 if count == 5000
  if line != "SKIP"
    `echo "100 100 100 100\n1\n#{line}" > problems#{i}/#{count}.txt`
  end
  count += 1
end