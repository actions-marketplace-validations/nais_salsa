package jvm

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nais/salsa/pkg/digest"
	"github.com/nais/salsa/pkg/scan"
)

func TestGradleDeps(t *testing.T) {
	got, _ := GradleDeps(gradleDependenciesOutput)
	want := []scan.Dependency{
		dep("ch.qos.logback:logback-classic", "1.2.9"),
		dep("org.apache.logging.log4j:log4j-core", "2.14.1"),
		dep("org.apache.logging.log4j:log4j-api", "2.14.2"),
		dep("org.wso2.orbit.org.zapache.commons:commons-collections4", "4.4.wso2v1"),
		dep("com.nimbusds:oauth2-oidc-sdk", "9.17"),
		dep("com.github.stephenc.jcip:jcip-annotations", "1.0-1"),
		dep("junit:junit", "3.8.1"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func dep(coordinates, version string) scan.Dependency {
	return scan.Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum: scan.CheckSum{
			Algorithm: "todo",
			Digest:    "todo",
		},
	}
}

const gradleDependenciesOutput = `testCompileClasspath - Compile classpath for source set 'test'.
+--- ch.qos.logback:logback-classic -> 1.2.9 (*)
+--- org.apache.logging.log4j:log4j-core:2.14.1
|    \--- org.apache.logging.log4j:log4j-api:2.14.2
+--- org.wso2.orbit.org.zapache.commons:commons-collections4:4.4.wso2v1
+--- com.nimbusds:oauth2-oidc-sdk:9.17
|    +--- com.github.stephenc.jcip:jcip-annotations:1.0-1
\--- junit:junit:3.8.1`

func TestGradleChecksums(t *testing.T) {
	buildToolMetadata := scan.CreateMetadata()
	err := GradleDepsAndSums(buildToolMetadata, []byte(gradleChecksumsOutput))
	if err != nil {
		fmt.Printf("%v", err)
	}
	want := map[string]scan.CheckSum{
		"ch.qos.logback:logback-classic": {Algorithm: digest.SHA256, Digest: "MzE2MGFlOTg4YWY4MmM4YmYzMDI0ZGRiZTAzNGE4MmRhOThiYjE4NmZkMDllNzZjNTA1NDNjNWI5ZGE1Y2M1ZQ=="},
		"ch.qos.logback:logback-core":    {Algorithm: digest.SHA256, Digest: "YmE1MWEzZmU1NjY5MWY5ZGQ3ZmU3NDJlNGE3M2MzYWI0YWFhYTMyMjAyYzczNDA5YmE1NmYxODY4NzM5OWEwOA=="},
	}

	if !reflect.DeepEqual(buildToolMetadata.Checksums, want) {
		t.Errorf("got %v, wanted %v", buildToolMetadata.Checksums, want)
	}
}

const gradleChecksumsOutput = `testCompileClasspath - Compile classpath for source set 'test'.
<?xml version="1.0" encoding="UTF-8"?>
<verification-metadata xmlns="https://schema.gradle.org/dependency-verification" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="https://schema.gradle.org/dependency-verification https://schema.gradle.org/dependency-verification/dependency-verification-1.1.xsd">
   <configuration>
      <verify-metadata>true</verify-metadata>
      <verify-signatures>false</verify-signatures>
   </configuration>
   <components>
      <component group="ch.qos.logback" name="logback-classic" version="1.2.10">
         <artifact name="logback-classic-1.2.10.jar">
            <sha256 value="3160ae988af82c8bf3024ddbe034a82da98bb186fd09e76c50543c5b9da5cc5e" origin="Generated by Gradle"/>
         </artifact>
         <artifact name="logback-classic-1.2.10.pom">
            <sha256 value="4accd4aa1c65d1acc70e6c391ab1033f5276b6e08628f1699d429aa4c3d5b884" origin="Generated by Gradle"/>
         </artifact>
      </component>
      <component group="ch.qos.logback" name="logback-core" version="1.2.10">
         <artifact name="logback-core-1.2.10.jar">
            <sha256 value="ba51a3fe56691f9dd7fe742e4a73c3ab4aaaa32202c73409ba56f18687399a08" origin="Generated by Gradle"/>
         </artifact>
         <artifact name="logback-core-1.2.10.pom">
            <sha256 value="1d32cf679d8780d7c9c915e6cc5db431d5b40c2aedc58a84a17cd023e1f0d09c" origin="Generated by Gradle"/>
         </artifact>
      </component>
   </components>
</verification-metadata>`
