problems = []
Dir.glob("problems/*").each do |f|
  /problems\/(\d+).txt/ =~ f
  problems << $1.to_i
end

num = problems.size
i = 1
c = 0
problems.sort.each do |f|
  command = "./Puzzle < problems/#{f}.txt | tee answers/#{f}.txt"
  result = `#{command}`
  c += 1 if result != "\n"
  
  print "#{f} [#{i}/#{num}] +#{c} "
  print command
  print " - "
  print result
  i += 1
end
