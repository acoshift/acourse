build: clean
	node build/build.js

dev:
	node build/dev-server.js

clean:
	rm -rf dist

prepare:
	npm install
	npm update

database:
	firebase deploy --only database

deploy: prepare build
	firebase deploy
