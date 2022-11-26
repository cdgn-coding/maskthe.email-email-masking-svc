deploy-minikube:
	minikube image build -t email-masking-svc:v1 .
	cd deployment && pulumi up -y