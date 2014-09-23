def symbolize_keys_deep!(h)
  h.keys.each do |k|
    ks    = k.to_sym
    h[ks] = h.delete k
    symbolize_keys_deep! h[ks] if h[ks].kind_of? Hash
  end
  h
end

