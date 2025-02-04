---
title: How to expose CloudQuery with PostGraphile
tag: tutorial
date: 2022/06/16
description: Search everything you have in the cloud with GraphQL.
author: yevgenypats
---

import { BlogHeader } from "../../components/BlogHeader"

<BlogHeader/>


In this blog post, we will walk you through how to setup [CloudQuery](https://github.com/cloudquery/cloudquery) to build your cloud asset inventory in PostgreSQL and build a GraphQL API query layer with [PostGraphile](https://github.com/graphile/postgraphile) on top of it. this can be used to build different use cases on from search to security, cost and infrastructure automation.


## **General Architecture**

- **ETL (Extract-Transform-Load) ingestion layer:** [CloudQuery](https://github.com/cloudquery/cloudquery)
- **Datastore:** PostgreSQL
- **API Access Layer:** [PostGraphile](https://github.com/graphile/postgraphile) and [GraphiQL](https://github.com/graphql/graphiql)

## What You Will Get

- **Raw SQL access** to all your cloud asset inventory to create views or explore any questions or connections between resources.
- **Multi-Cloud Asset Inventory:** Ingest configuration from all your clouds to a single datastore with a unified structure.
- **GraphQL Endpoint** to access and query all your cloud configurations.

## Walkthrough

### Step 1: **Install or Deploy CloudQuery (fetch)**

If it’s your first time using CloudQuery we suggest you first run it locally to get familiar with the tool. Take a look at our [Getting Started with AWS Guide](https://docs.cloudquery.io/docs/getting-started/getting-started-with-aws) or [GCP](https://docs.cloudquery.io/docs/getting-started/getting-started-with-gcp), [Azure](https://docs.cloudquery.io/docs/getting-started/getting-started-with-azure) .

If you are already familiar with CloudQuery, take a look at how to deploy it to AWS on RDS Aurora and EKS at [github.com/cloudquery/terraform-aws-cloudquery](https://github.com/cloudquery/terraform-aws-cloudquery) .

### Step 2: Install PostGraphile

For full details, check out the [PostGraphile](https://www.graphile.org/postgraphile/quick-start-guide/) docs. If you are running locally will need **Node.js** and you can install **PostGraphile** globally via `npm i -g postgraphile` or (`brew install PostGraphile`)

To run PostGraphile locally all you need to do is the following (adjust the Postgres URL accordingly):

```bash
postgraphile -c "postgres://postgres:pass@localhost:5432/postgres" --enhance-graphiql --skip-plugins graphile-build:NodePlugin --simple-collections only -p 6060
```

### Step 3: Query and Profit!

That’s it! The output of a successful run is presented below:

```bash
PostGraphile v4.12.9 server listening on port 6060 🚀

  ‣ GraphQL API:         http://localhost:6060/graphql
  ‣ GraphiQL GUI/IDE:    http://localhost:6060/graphiql
  ‣ Postgres connection: postgres://postgres:[SECRET]@localhost/postgres
  ‣ Postgres schema(s):  public
  ‣ Documentation:       https://graphile.org/postgraphile/introduction/
  ‣ Node.js version:     v18.3.0 on darwin arm64
  ‣ Join PostHog in supporting PostGraphile development: https://graphile.org/sponsor/

* * *
```

Open the browser with the `http://localhost:6060/graphiql` endpoint to see the GraphiQL UI where you can compose any query you want interactively:

![](/images/blog/cloudquery-postgraphile/step3.png)

### Step 4: Create New Views

By default PostGraphile exposes all tables and relationships of the existing tables but let’s say you want to create a new view. All you need to do is to create a new view and PostGraphile will automatically generate the model for that. For example, check out this [blog](https://www.cloudquery.io/blog/aws-resources-view) on how to create a unified AWS resource [view](https://github.com/cloudquery/cq-provider-aws/tree/main/views) (or GCP [View](https://github.com/cloudquery/cq-provider-gcp/tree/main/views)). And just like that you can now query and search all your resources by `arn`, `tags`, `name` with GraphQL!

![](/images/blog/cloudquery-postgraphile/step4.png)

### Step 5: Deploying in production

If you want to expose PostGraphile publicly please see [PostGraphile](https://www.graphile.org/postgraphile/security/) Security or expose it privately and use either a bastion host or something like [Tailscale on Kubernetes](https://tailscale.com/kb/1185/kubernetes/) together with our [helm charts.](https://github.com/cloudquery/helm-charts)

## Summary

In this post we showed you how to build an open-source cloud asset inventory with CloudQuery as the ETL (Extract-Transform-Load) / data-ingestion layer and [PostGraphile](https://www.graphile.org/) as the API layer to expose the data for your internal team/users or any other downstream processing in the most convenient/preferred way.
