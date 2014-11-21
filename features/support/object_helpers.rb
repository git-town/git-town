def symbolize_keys_deep! hash
  hash.keys.each do |key|
    symbol_key = key.to_s.gsub(' ', '_').to_sym
    hash[symbol_key] = (value = hash.delete key)
    symbolize_keys_deep! value if value.is_a? Hash
  end
  hash
end
