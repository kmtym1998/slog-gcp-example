# .env をロード
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: deploy
deploy: # デプロイ
	make build-for-deploy
	make cloud-run

.PHONY: build-for-deploy
build-for-deploy: # イメージのビルド
	docker buildx build --platform linux/amd64 -t asia.gcr.io/$(PROJECT_ID)/slog-example .

.PHONY: cloud-run
cloud-run: # イメージを GCR に push して Cloud Run にデプロイ
	docker push asia.gcr.io/$(PROJECT_ID)/slog-example
	gcloud run deploy slog-example \
		--region=asia-northeast1 \
		--allow-unauthenticated \
		--image=asia.gcr.io/$(PROJECT_ID)/slog-example:latest \
		--set-env-vars=PROJECT_ID=$(PROJECT_ID) \
		--min-instances=0 \
		--max-instances=1 \
		--memory=128Mi \
		--cpu=1

.PHONY: run-container
run-container: # ローカルでコンテナを動かす
	docker build -t slog-example-local .
	docker run slog-example-local
