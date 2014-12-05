# Extend the hash class
class Hash
  # Converts all keys to symbols
  def symbolize_keys_deep!
    keys.each do |key|
      symbol_key = key.to_s.gsub(' ', '_').to_sym
      self[symbol_key] = (value = delete key)
      value.symbolize_keys_deep! if value.is_a? Hash
    end
  end


  # Replaces any blank values with the defaults
  def default_blank! defaults
    defaults.each_pair do |key, value|
      self[key] = value if self[key].blank?
    end
  end

  # Creates a subhash with just the following keys
  def subhash *keys
    select { |key| keys.include?(key) }
  end
end
