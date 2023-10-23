cert:
	cd cert; chmod +x gen.sh; sudo -S ./gen.sh; sudo chmod -R 777 *.pem; cd ..

.PHONY: gen clean server client test cert