# Moments
![Current master build status](https://travis-ci.org/valerykalashnikov/moments.svg?branch=master)


Embeddable timeseries storage to track moments for n previous minutes/seconds/milliseconds

## Why would I Need It?

In any place you need to implement counter for n previous seconds. For example to display requests per previous 60 seconds to provide section with addition information in /healhcheck endpoints.

## How to use it?

To define counter without using file to backup:

~~~go
  // initialize counter
  counter = moments.NewMomentsCounter(1 * time.Minute)

  // track moment
  counter.Track()
  counter.Track()

  // display value
  val := counter.Count()
  fmt.Println(val)
~~~

To define counter using file to backup:

~~~go

  // open file with backup
  f, err := os.OpenFile("/path/to/backup", os.O_RDWR, os.ModeAppend)
  if err != nil {
    log.Fatalf("Unable to open file to backup moments, %s", err)
  }
  defer f.Close()

  // initialize counter
  counter, err := moments.NewMomentsCounterFrom(f)
  if err != nil {
    // handle error
  }

  counter.Track()
  counter.Track()

  // save value to backup file
  err := counter.Save(f)

  if err != nil {
    //handle error
  }
~~~

Here in [example](https://github.com/valerykalashnikov/moments/tree/master/example) folder there is a simple http application with unit/e2e test to show the approach.
