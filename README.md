# workout-helper

Small tool to help with workout timing, and guided sections with a simple robo-coach.

## Back story

In spring, in the the lockdown, I did workout at home based on the workout scripts that my kettlebell instructor gave me.

The workouts followed similar patterns, with little variance.
Warmup at the start and stretch at the end were part of every session.
Between those the sections had different approximate times.
Once per week we had a tabata session with isometric movements.

Now that the pandemic lockdown has arrived again, I come prepared for the second act.

I'm happy with the result, from now on I do not need to check the time again, and a calm voice guides me through these workouts.

## Getting started

Currently I do not ship binaries, or docker image, so in order to run it you need to have golang installed, I recommend the latest version.

The base workout consists of 10 minutes of warmup, 50 minutes of a general task, 10 minutes of stretching.

```bash
go run .
```

The code can accept one string parameter. Tabata is marked with `t`, and general sections are marked with `s\d+`.
They have to be separated by a space character. Order counts.

e.g: `'t s10 s20'` means:

1. Warmup
1. A Tabata session
1. 10 minutes for one task
1. 20 minutes for another task
1. Stretching

```bash
go run . 't s10 s20'
```
