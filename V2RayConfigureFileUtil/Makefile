pb:
		@echo "pb Start"
		cd encoding && make pb

mkseg:
	-mkdir seg

mkout:
		-mkdir out

out/%.qr.png:seg/%.seg.libv2ray.pb
		./V2RayConfigureFileUtil -t QR -i $^ -o $@

out/%.az.png:seg/%.seg.libv2ray.pb
				./V2RayConfigureFileUtil -t AZ -i $^ -o $@

seg: mkseg
	./V2RayConfigureFileUtil -t seg -i inp.json -o seg/

srcfiles := $(shell echo seg/*.seg.libv2ray.pb)
destfiles := $(patsubst seg/%.seg.libv2ray.pb,out/%.qr.png,$(srcfiles))
destfilesaz := $(patsubst seg/%.seg.libv2ray.pb,out/%.az.png,$(srcfiles))

dispo:
	@echo "$(srcfiles)"

convqr: mkout dispo $(destfiles)
	@echo "Done"

convaz: mkout dispo $(destfilesaz)
		@echo "Done"

genqr: clean seg
		$(MAKE) convqr
		#$(MAKE) convaz
		@echo "Done"

clean:
		-rm -R out seg

ShippedBinary:
	cd Convert; $(MAKE) shippedBinary

all: ShippedBinary pb
	@echo OK!
