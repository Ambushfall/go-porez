#!/usr/bin/env bash

package=$1
all=$2
if [[ -z "$package" ]]; then
	echo "usage: $0 <package-name>"
	exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64")

if [[ -n "$all" && "$all" == "ALL" ]]; then

	testplatforms=$(go tool dist list)
	platforms=(${testplatforms// //})

	delete=("android/386" "android/amd64" "android/arm" "android/arm64" "ios/amd64" "ios/arm64")

	for target in "${delete[@]}"; do
		for i in "${!platforms[@]}"; do
			if [[ ${platforms[i]} = $target ]]; then
				unset 'platforms[i]'
			fi
		done
	done

	platforms=("${platforms[@]}")
fi

for platform in "${platforms[@]}"; do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	echo "GOOS:$GOOS GOARCH:$GOARCH outname: $output_name"
	env GOOS=$GOOS GOARCH=$GOARCH go build -o $GOOS/$GOARCH/$output_name .
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done
