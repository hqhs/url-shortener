Classic service example: URL shortener
===

Average reading time: 10-15 minutes.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Classic service example: URL shortener](#classic-service-example-url-shortener)
    - [Abstract](#abstract)
    - [Features](#features)
    - [Database choice](#database-choice)
    - [URL shortener algorithm choice](#url-shortener-algorithm-choice)
    - [Data partitioning](#data-partitioning)
    - [Other possible improvements](#other-possible-improvements)
    - [Additional services](#additional-services)
    - [Actual deployment configuration with Kubernetes](#actual-deployment-configuration-with-kubernetes)
    - [Contributing](#contributing)
        - [Run in docker](#run-in-docker)
        - [Run without docker](#run-without-docker)

<!-- markdown-toc end -->

## Abstract

This is thought experiment about engineering url-shortener service with high
load capabilities. Any corrections are welcome!

## Database choice

  Let's say average long URL size is 2kb. Additionally I want to record
statistics, such as Unix timestamp and IP request address for each click. Keep
in mind that I may add more statistics in future. 10 bytes for timestamp and 16
bytes for IP (with IPv6 support) sounds reasonable. Statistics is not
implemented at all (yet), but in the world of data driven advertising loosing
opportunity to collect useful info is not very realistic.

bit.ly statistics from ([source]( http://highscalability.com/blog/2014/7/14/bitly-lessons-learned-building-a-distributed-system-that-han.html )): 
  - 6 * 10^9 clicks/month
  - 6 * 10^8 shortens/month,
  - avarage click/link is 10, but maximum click per link could be anything.
 
 6 * 10^9 * 26 bytes + 6 * 10^8 * 2 kb = 1.35 terabytes/month
 
 I need to choose database wisely. Take into account what data has high
read/write ratio, never change, not structured (even with users structure would
be slightly structured), and record size is small (around 2kb).
 
 I don't need relational database -- data is too simple, not column store --
there's no need for searching/sorting, not graph db -- there's no graph, not
document store -- this is not the fastest available option. What I need is
key/value store, which is reliable out of the box, easy to scale, and easy to
get up and running.

 So, for bit.ly comparable traffic I can use one of those: AerospikeDB, Apache
Cassandra, Riak, DynamoDB (those are most popular, there're other choices). 

 Remember, what statistics above taken from startup which is **primary** product
is URL shortening and analytics service. Just copy-pasting their architecture
(and choosing hadoop) solutions doesn't fit because I can't **expect** such
traffic. And I'm not competent enough for building infrastructure capable to
withstand **real** high load traffic (though someday I will), so for purpose of
making working demo I'm going with redis cluster and implement database
connection as golang interface, which allows us to quickly change database then
needed.

## URL shortener algorithm choice

- **Option 0: Simple counter stored in database**:
Store counter in database. To enforce shorten URLs be longer then N chars,
counter should start with 64^N, since service use base64 url-safe encoding.
For new URL get current counter value, encode it, and use as key. Increment
after usage. 

Issue: It's easy to iterate over all URLs.
- **Option 1: Hash based algorithm with counter**: 
Hash requested URL, base64 encode it, take first N chars, check what received 
key is unique in database. To enforce different keys for same URLs hash with
counter technique from previous variant.

Issue: It's slows down over time. Let's calculate degradation speed.
Given N=6 there are ~6^64 or ~10^50 possible keys. With 6 * 10^8 request/month
(bit.ly load rate from above), that's 7.2 * 10^9 per year and ~10^11 per 15
years. After 15 years 10^11 / 10^50 = 1/10^39 of all possible keys will be used.
Practically speaking, every new key is random and has 1/10^39 probability to be 
already used. So, on avarage, every 1/10^39 request makes additional database
request after 15 years of bit.ly load rate.

- **Option 2: Zookeper**:
Use separate (let's call it zookeper) service which pre-generates keys for
service nodes. Then node bootstraps or run out of keys it requests new set of
unique keys from zookeper and use it.

Issue: Zookeper is a single point of failure and adds unnecessary entities.

First option has good price/perfomance ration. I'm going with it.

## Data partitioning

  Data partitioning is heavily based on database of your choice, but with 
"practically speaking" random URL shortening algorithm service could 
split urls between master nodes (if any) based on their shorten versions.

## Other possible improvements

- Caching. We can start with 20% caching of daily traffic and adjust that value then
needed. Least Recently Used (LRU) is reasonable policy for our system.
- Load balancing. Initially simple Round Robin approach would suffice. After
noticeable traffic growth more complex algorithm with back pressure support is
needed.
- Throttling and graceful degradation. 
- Telemetry.
- Health checking.

## Additional services

  It's a good idea to move data gathering to separate service with column based
database (such as ClickHouse). For this particular feature we need to add a
method to distinguish users from each other. For example, generate unique
cookie (or allow users to register themselves) for each user, build profile with
their interests based on their requested URLs, and show advertising based on
this profile (or just sell data). Simplest way to find correlation between URL
and interests is to parse destination page's SEO metatags. Other way to monetize
(bit.ly uses it) is analytics selling for registered users (such as click rate etc.)

## Actual deployment configuration with Kubernetes

## Contributing

### Run in docker

### Run without docker

