dir = "answers/"

for i in 1..5000
  begin
    f = open("#{dir}#{i}.txt")
    answer = f.gets
    puts answer
  rescue
    puts
  end
end
