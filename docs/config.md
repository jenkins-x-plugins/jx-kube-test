---
title: API Documentation
linktitle: API Documentation
description: Reference of the jx-promote configuration
weight: 10
---
<p>Packages:</p>
<ul>
<li>
<a href="#kubetest.jenkins-x.io%2fv1alpha1">kubetest.jenkins-x.io/v1alpha1</a>
</li>
</ul>
<h2 id="kubetest.jenkins-x.io/v1alpha1">kubetest.jenkins-x.io/v1alpha1</h2>
<p>
<p>Package v1alpha1 is the v1alpha1 version of the API.</p>
</p>
Resource Types:
<ul><li>
<a href="#kubetest.jenkins-x.io/v1alpha1.KubeTest">KubeTest</a>
</li></ul>
<h3 id="kubetest.jenkins-x.io/v1alpha1.KubeTest">KubeTest
</h3>
<p>
<p>KubeTest represents the configuration of kube test</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
kubetest.jenkins-x.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>KubeTest</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.KubeTestSpec">
KubeTestSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Spec holds the desired state of the KubeTest from the client</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>rules</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Rule">
[]Rule
</a>
</em>
</td>
<td>
<p>Rules the rules to apply</p>
</td>
</tr>
<tr>
<td>
<code>outputDir</code></br>
<em>
string
</em>
</td>
<td>
<p>OutputDir the output directory to store the reports</p>
</td>
</tr>
<tr>
<td>
<code>format</code></br>
<em>
string
</em>
</td>
<td>
<p>Format the output format</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="kubetest.jenkins-x.io/v1alpha1.Charts">Charts
</h3>
<p>
(<em>Appears on:</em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Rule">Rule</a>)
</p>
<p>
<p>Charts the charts to template and validate</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>dir</code></br>
<em>
string
</em>
</td>
<td>
<p>Dir the directory containing a helm chart or the source to recurse through if recursive is enabled</p>
</td>
</tr>
<tr>
<td>
<code>recurse</code></br>
<em>
bool
</em>
</td>
<td>
<p>Recurse if enabled recurse through the directory to find any Chart.yaml files</p>
</td>
</tr>
</tbody>
</table>
<h3 id="kubetest.jenkins-x.io/v1alpha1.KubeTestSpec">KubeTestSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#kubetest.jenkins-x.io/v1alpha1.KubeTest">KubeTest</a>)
</p>
<p>
<p>KubeTestSpec defines the configuration of kube test</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>rules</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Rule">
[]Rule
</a>
</em>
</td>
<td>
<p>Rules the rules to apply</p>
</td>
</tr>
<tr>
<td>
<code>outputDir</code></br>
<em>
string
</em>
</td>
<td>
<p>OutputDir the output directory to store the reports</p>
</td>
</tr>
<tr>
<td>
<code>format</code></br>
<em>
string
</em>
</td>
<td>
<p>Format the output format</p>
</td>
</tr>
</tbody>
</table>
<h3 id="kubetest.jenkins-x.io/v1alpha1.Rule">Rule
</h3>
<p>
(<em>Appears on:</em>
<a href="#kubetest.jenkins-x.io/v1alpha1.KubeTestSpec">KubeTestSpec</a>)
</p>
<p>
<p>Rule the rules to apply</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>resources</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Source">
Source
</a>
</em>
</td>
<td>
<p>Resources the kubernetes resource dir to look for resources to verify</p>
</td>
</tr>
<tr>
<td>
<code>charts</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Charts">
Charts
</a>
</em>
</td>
<td>
<p>Charts the charts to evaluate</p>
</td>
</tr>
<tr>
<td>
<code>tests</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Tests">
Tests
</a>
</em>
</td>
<td>
<p>Tests the tests to perform</p>
</td>
</tr>
</tbody>
</table>
<h3 id="kubetest.jenkins-x.io/v1alpha1.Source">Source
</h3>
<p>
(<em>Appears on:</em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Rule">Rule</a>)
</p>
<p>
<p>Source the location of kubernetes resources to validate</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>dir</code></br>
<em>
string
</em>
</td>
<td>
<p>Dir the directory containing the kubernetes resources</p>
</td>
</tr>
</tbody>
</table>
<h3 id="kubetest.jenkins-x.io/v1alpha1.Test">Test
</h3>
<p>
(<em>Appears on:</em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Tests">Tests</a>)
</p>
<p>
<p>Test a kind of test</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>version</code></br>
<em>
string
</em>
</td>
<td>
<p>Version optional override of the version to use</p>
</td>
</tr>
<tr>
<td>
<code>args</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Args optional additional comand line arguments to pass to the test</p>
</td>
</tr>
</tbody>
</table>
<h3 id="kubetest.jenkins-x.io/v1alpha1.Tests">Tests
</h3>
<p>
(<em>Appears on:</em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Rule">Rule</a>)
</p>
<p>
<p>Tests the tests to run on the resources</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conftest</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Test">
Test
</a>
</em>
</td>
<td>
<p>Conftest enables conftest tests</p>
</td>
</tr>
<tr>
<td>
<code>kubescore</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Test">
Test
</a>
</em>
</td>
<td>
<p>Kubescore enables kube-score based tests</p>
</td>
</tr>
<tr>
<td>
<code>kubeval</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Test">
Test
</a>
</em>
</td>
<td>
<p>Kubeval enables kubeval tests</p>
</td>
</tr>
<tr>
<td>
<code>polaris</code></br>
<em>
<a href="#kubetest.jenkins-x.io/v1alpha1.Test">
Test
</a>
</em>
</td>
<td>
<p>Polaris enables polaris tests</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>bebe38e</code>.
</em></p>
