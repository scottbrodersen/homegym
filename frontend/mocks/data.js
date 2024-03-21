export const fetchedEvents = (pageSize) => {
  events = [];
  for (let i; i < pageSize; i++) {
    events.push({
      id: `eventID${i}`,
      activityID: `activity${randomIntFromInterval(1, 3)}`,
      date: Math.floor(Date.now() / 1000) - pageSize * 1000 + i * 100,
      mood: randomIntFromInterval(1, 5),
      motivation: randomIntFromInterval(1, 5),
      energy: randomIntFromInterval(1, 5),
      overall: randomIntFromInterval(1, 5),
      notes: `test note ${i}`,
      exercises: null,
    });
  }
};

export const randomIntFromInterval = (min, max) => {
  // min and max included
  return Math.floor(Math.random() * (max - min + 1) + min);
};

export const fetchedTestActivities = [
  {
    id: 'activity1',
    name: 'test activity 1',
    exercises: ['exercise1', 'exercise2'],
  },
  {
    id: 'activity2',
    name: 'test activity 2',
    exercises: null,
  },
  {
    id: 'activity3',
    name: 'test activity 3',
    exercises: null,
  },
];

export const fetchedExercises = [
  {
    name: 'test exercise 1',
    id: 'exercise1',
    intensityType: 'weight',
    volumeType: 'count',
    volumeConstraint: 2,
    composition: null,
    basis: '',
  },
  {
    name: 'test exercise 2',
    id: 'exercise2',
    intensityType: 'weight',
    volumeType: 'count',
    volumeConstraint: 1,
    composition: null,
    basis: '',
  },
];

export const testIntensity = (type) => {
  let intensity;
  if (type === 'weight') {
    intensity = (randomIntFromInterval(50, 200) / 3).toFixed(1);
  }
  return intensity;
};

export const testVolume = () => {
  const sets = [];
  for (let set = 0; set < randomIntFromInterval(2, 6); set++) {
    const reps = [];
    for (let rep = 0; rep < randomIntFromInterval(2, 8); rep++) {
      reps.push(randomIntFromInterval(0, 1));
    }
    sets.push(reps);
  }
  return sets;
};

export const testExerciseInstance = (index) => {
  const parts = [];
  for (let i = 0; i < randomIntFromInterval(1, 4); i++) {
    parts.push({
      intensity: testIntensity(),
      volume: testVolume(),
    });
  }
  return {
    typeID:
      fetchedExercises[randomIntFromInterval(0, fetchedExercises.length)].id,
    index: index,
    parts: parts,
  };
};

export const fetchedEventExercises = () => {
  instances = new Map();
  for (let i = 0; i < randomIntFromInterval(1, 5); i++) {
    instances.set(i, testExerciseInstance(i));
  }
  return instances;
};

export const testProgram = () => {
  return {
    id: 'test-program-id',
    activityID: 'test-activity-id',
    title: 'test program title',
    blocks: [
      {
        title: 'test block',
        microCycles: [
          {
            title: 'test Microcycle',
            workouts: [
              {
                title: 'test workout',
                segments: [
                  { title: 'test segment', prescription: 'test prescription' },
                ],
              },
            ],
          },
        ],
      },
    ],
  };
};
