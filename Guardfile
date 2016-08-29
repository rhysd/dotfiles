def run_test(file)
  system "go test -v ./#{file} #{Dir['./dotfiles-command/*.go'].reject{|p| p.end_with? '_test.go'}.join(' ')}"
end

guard :shell do
  watch /\.go$/ do |m|
    puts "\033[93m#{Time.now}: #{File.basename m[0]}\033[0m"
    case m[0]
    when /_test\.go/
      run_test m[0]
    else
      test = m[0][0..-4] + "_test.go" # foo.go -> foo_test.go
      if File.exists? test
        run_test(test)
      else
        system "go build"
      end
    end
  end
end
