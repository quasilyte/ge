.PHONY: update-ebitengine-input
update-ebitengine-input:
	rm -rf ./input
	git clone --depth=1 https://github.com/quasilyte/ebitengine-input.git input
	rm -rf ./input/.git
	rm ./input/go.mod && rm ./input/go.sum
	rm ./input/math_nodeps.go
	tail -n +3 ./input/math_gmath.go > tmp && cat tmp > ./input/math_gmath.go && rm tmp
