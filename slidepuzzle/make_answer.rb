answer = "answers.txt"

open(answer).each do |l|
  new_line = gets.chop
  
  if new_line != ""
    puts new_line
  else
    puts l
  end
end