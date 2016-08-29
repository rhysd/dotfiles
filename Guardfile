def run_test(file)
end

guard :shell do
  watch /\.go$/ do |m|
    puts "\033[93m#{Time.now}: #{File.basename m[0]}\033[0m"
    case m[0]
    when /_test\.go$/
      system "go test -v ./#{m[0]} #{Dir['./dotfiles-command/*.go'].reject{|p| p.end_with? '_test.go'}.join(' ')}"
    else
      system "go build"
    end
  end
end
