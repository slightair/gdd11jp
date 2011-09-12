u = d = l = r = n = 0
while line = gets
  line.chop!
  n += 1 if line.size > 0
  line.bytes.each do |b|
    case b
    when ?U
      u += 1
    when ?D
      d += 1
    when ?L
      l += 1
    when ?R
      r += 1
    end
  end
end

puts "answers:#{n}/5000(#{"%.3f"%(n/5000.0*100)}%) [LX:#{l}(#{"%.3f"%(l/72187.0*100)}%) RX:#{r}(#{"%.3f"%(r/81749.0*100)}%) UX:#{u}(#{"%.3f"%(u/72303.0*100)}%) DX:#{d}(#{"%.3f"%(d/81778.0*100)}%)]"