# gauther

Basic Google OAuth2 implementation with [Firestore](https://cloud.google.com/firestore/) persistence at global scale

## Demo

https://gauther.default.knative.tech/

## Setup

> NOTE: Work in progress

For now, take a look at the [Makefile](https://github.com/mchmarny/gauther/blob/master/Makefile) for help with the major configuration commands. 


### Knative URL

To avoid the kind of chicken and an egg situation we are going to first define the `URL` that your application will have when you publish it on Knative. Knative uses convention to build serving URL by combining the deployment name (e.g. `gauther`), namespace name (e.g. `default`), and the pre-configured domain name (e.g. `knative.tech`). The resulting URL, assuming you already configured SSL, should look something like this:

```shell
https://gauther.default.knative.tech
```

### Google OAuth Credentials

In your Google Cloud Platform (GCP) project console navigate to the Credentials section. You can use the search bar, just type `Credentials` and select the option with "API & Services". To create new OAuth credentials:

* Click “Create credentials” and select “OAuth client ID”
* Select "Web application"
* Add authorized redirect URL at the bottom using the fully qualified domain we defined above and appending the `callback` path:
 * `https://gauther.default.knative.tech/auth/callback`
* Click create and copy both `client id` and `client secret`
* CLICK `OK` to save

> You will also have to verify the domain ownership. More on that [here](https://support.google.com/cloud/answer/6158849?hl=en#authorized-domains)

### Google Cloud Firestore

If you haven't used yet Firestore on GCP, you will have to enable it. You can find instructions on how to do it [here](https://firebase.google.com/docs/firestore/quickstart)

### App Configuration

