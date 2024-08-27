# My research on data

Assuming that description represents information about the file accurately and in full,
this is the query to retrieve the list of machine-readable file URLs corresponding to Anthem's PPO in New York state:
```clickhouse
SELECT description, location
FROM urls
WHERE (lower(description) LIKE '%new york%' or lower(description) LIKE '%ny%')
  AND lower(description) LIKE '%ppo%';
```

Let's count them:
```clickhouse
SELECT count(*)
FROM urls
WHERE (lower(description) LIKE '%new york%' or lower(description) LIKE '%ny%')
  AND lower(description) LIKE '%ppo%';

-- Result: 542835
```

However, the file location URLs are looking like they are built with pattern `"anthem" + ("bc" or "bcbs") + state code`.
Except some strange links like this (none of them are relevant to NY):
```clickhouse
SELECT description, count(location)
FROM urls
WHERE lower(location) NOT LIKE '%bcbs%'
GROUP BY description;

-- Result:
-- In-Network Negotiated Rates Files,1077809
-- Health Service Coalition of Nevada (HSC) Rates,1
-- Dental Vision,204512
-- Beacon Behavioral Health Rates Files,7
```

I also checked out the Anthem EIN lookup page, where I input Highmark in the company name field 
and look up for the MRF URLs. It contained list of links, some of them had "NY" in the name,
but when I copied the link address, it was pointing to the file with no such word as "NY", "New York" or "PPO",
it was actually the opposite - the links contained "nv" for the state code. In conclusion, I can't say that
I can rely on the "location" field to determine if it is representing NY or not.

By the way, I didn't find and of "NY" as state code part of location:
```clickhouse
SELECT count(*)
FROM urls
WHERE match(lower(location), 'https://anthem[a-z]*ny\..*');

-- Result:
-- 0

-- There are 899k of NV locations:
SELECT count(*)
FROM urls
WHERE match(lower(location), 'https://anthem[a-z]*nv\..*');

-- Result:
-- 899362
```

"Is Highmark the same as Anthem?". There is no rows with "description" containing both "NY" and "PPO",
but not containing "Highmark".   

```clickhouse
SELECT count(*)
FROM urls
WHERE (lower(description) LIKE '%new york%' or lower(description) LIKE '%ny%')
  AND lower(description) NOT LIKE '%highmark%';

-- Result:
-- 0
```

However, it is not the case for all the other rows:

```clickhouse
SELECT count(*)
FROM urls
WHERE lower(description) NOT LIKE '%highmark%';

-- Result:
-- 36796423
```

So, Highmark is not the same as Anthem, but might be the exclusive representative of it in NY

In total, I could not find a correlation between "location" field and geographical location.
My final query of finding New York PPOs would be this:
```clickhouse
SELECT description, location
FROM urls
WHERE (lower(description) LIKE '%new york%' or lower(description) LIKE '%ny%')
  AND lower(description) LIKE '%ppo%';
```
