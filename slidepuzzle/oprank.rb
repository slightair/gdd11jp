answers = []
count = 1
while line = gets
  line.chop!
  answers << "%04d - %04d - %s" % [count, line.size, line]
  count += 1
end

answers.sort! do |a, b|
  a.size <=> b.size
end

answers.each do |a|
 puts a
end