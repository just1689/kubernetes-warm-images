helm package warm-images
mv *.tgz ../charts/warm-images/
helm repo index ../charts/warm-images --url https://storage.googleapis.com/captains-charts/warm-images