#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import requests
import json
import yaml
from tempfile import TemporaryDirectory
from subprocess import Popen


class ReleaseDataMissing(Exception):
    def __init__(self, dist, version):
        self.dist = dist
        self.version = version
        super().__init__(f"Configuration for released {dist} version {version} is missing")


class ReleaseDataChanged(Exception):
    def __init__(self, dist, version):
        self.dist = dist
        self.version = version
        super().__init__(f"Images for released {dist} version {version} have changed")


class UnexpectedHelmChart(Exception):
    def __init__(self, chart, version, repo, url):
        self.chart = chart
        self.version = version
        self.repo = repo
        self.url = url
        super().__init__(f"Unexpected chart URL for {repo}/{chart}@{version}: {url}")


def main(*releases):
    dev_data = json.load(open('data/data.json'))

    for release in releases:
        release_data_url = f"https://releases.rancher.com/kontainer-driver-metadata/{release}/data.json"
        response = requests.get(release_data_url)
        response.raise_for_status()
        release_data = response.json()

        check_rke(dev_data, release_data)
        check_k3s(dev_data, release_data)
        check_rke2(dev_data, release_data)


def check_rke(dev_data, release_data):
    dist = 'Rancher k8s'
    dev_images = dev_data.get('K8sVersionRKESystemImages', {})
    release_images = release_data.get('K8sVersionRKESystemImages', {})

    check_versions(dev_images, release_images, dist)


def check_k3s(dev_data, release_data):
    dist = 'k3s'
    dev_versions = get_releases(dev_data, dist)
    release_versions = get_releases(dev_data, dist)

    check_versions(dev_versions, release_versions, dist)


def check_rke2(dev_data, release_data):
    dist = 'rke2'
    dev_versions = get_releases(dev_data, dist)
    release_versions = get_releases(dev_data, dist)

    check_versions(dev_versions, release_versions, dist)
    check_rke2_charts(dev_versions)


def check_rke2_charts(versions):
    for version, config in versions.items():
        charts = config.get('charts', {})
        if not charts:
            continue

        image = f"docker.io/rancher/rke2-runtime:{version.replace('+', '-')}"
        with TemporaryDirectory(version) as tmp:
            with Popen(['wharfie', image, tmp]) as wharfie:
                retcode = wharfie.wait()
                if retcode != 0:
                    print(f"Unable to extract rke2 runtime image {image}; skipping chart validation")
                    continue

            for chart, info in charts.items():
                print(f"Checking RKE2 {version} {info['repo']}/{chart}@{info['version']}")
                with open(f"{tmp}/charts/{chart}.yaml") as manifest:
                    manifest_yaml = yaml.safe_load(manifest)
                    chart_url = manifest_yaml.get('metadata', {}).get('annotations', {}).get('helm.cattle.io/chart-url')
                    expected_chart_tarball = f"{chart}-{info['version']}.tgz"
                    if expected_chart_tarball not in chart_url:
                        raise UnexpectedHelmChart(chart, info['version'], info['repo'], chart_url)


def get_releases(data, dist):
    return {r['version']: r for r in data.get(dist, {}).get('releases', [])}


def check_versions(dev, release, dist):
    for version, config in release.items():
        print(f"Checking {dist} {version}")
        if dist in ['k3s', 'rke2']:
            # k3s and rke2 releases may need to be fixed after the fact; just make sure we don't remove a released version
            if version not in dev:
                raise ReleaseDataMissing(dist, version)
        else:
            # Rancher kubernetes release images cannot be changed after the fact
            if dev[version] != config:
                raise ReleaseDataChanged(dist, version)


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <release> [<release> ...]")
        exit(1)

    main(*sys.argv[1:])
