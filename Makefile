deploy-minikube:
	minikube image build -t email-masking-svc .
	cd deployment && pulumi up -y