helm package warm-images
mv *.tgz ../charts/
helm repo index  ../charts/ --url https://storage.googleapis.com/captains-charts/
# TODO: Copy dir to cloud