problem = "problems.txt"

count = 0
open(problem).each do |l|
  if count < 2
    puts l
  else
    ans_line = gets.chop
    
    if ans_line != ""
      puts "SKIP"
    else
      puts l
    end
  end
  count += 1
end