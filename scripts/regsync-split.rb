#!/usr/bin/env ruby

require "json"
require "pathname"
require "yaml"
require "pp"

# Determine paths
script_dir = Pathname.new(__FILE__).realpath.dirname
repo_root = script_dir.parent
config_file_path = repo_root + "regsync.yaml"

pp script_dir
pp repo_root
pp config_file_path

unless config_file_path.exist?
  abort "Error: Config file not found at: #{config_file_path}"
end

# Load the main config
regsync = YAML.load_file(config_file_path)

# Count total tags
sum = regsync["sync"].sum do |sync|
  allow_tags = sync.dig("tags", "allow") || []
  allow_tags.count
end

puts "Total tags to consider: #{sum}"

# Split each sync block
regsync["sync"].each do |sync|
  single_sync_config = regsync.merge("sync" => [sync])

  target_dir = repo_root + "split-regsync" + sync["source"]
  target_dir.mkpath

  output_file = target_dir + "split-regsync.yaml"
  output_file.write(YAML.dump(single_sync_config))

  puts "Wrote split config: #{output_file}"
end

