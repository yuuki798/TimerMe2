'use client';
// Home.tsx

import React, { useState, useEffect } from 'react';
import {
  getTasks,
  createTask,
  updateTask,
  deleteTask,
  startTask,
  pauseTask,
  completeTask,
  resetTask, // 新增重置任务API调用
} from '@/utils/api';
import TaskList from '@/components/TaskList';

interface Task {
  id: number;
  name: string;
  duration: number;
  is_completed: boolean;
  start_time: string;
  status: string;
  total_time: number;
}

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [newTaskName, setNewTaskName] = useState('');
  const [newTotalTime, setNewTotalTime] = useState<number>(0);
  const [taskTimers, setTaskTimers] = useState<{
    [key: number]: NodeJS.Timeout;
  }>({});
  const [timeUnit, setTimeUnit] = useState('seconds');

  useEffect(() => {
    fetchTasks();
  }, []);

  useEffect(() => {
    tasks.forEach((task) => {
      if (task.status === 'started' && !taskTimers[task.id]) {
        const startTime = new Date(task.start_time).getTime();
        const currentTime = Date.now();
        const elapsedTime = Math.floor((currentTime - startTime) / 1000);
        startTaskTimer(task.id, elapsedTime, task.duration);
      }
    });
  }, [tasks]);

  const startTaskTimer = (
    id: number,
    initialElapsedTime: number = 0,
    initialDuration: number = 0,
  ) => {
    const startTime = Date.now() - initialElapsedTime * 1000;
    const timer = setInterval(() => {
      setTasks((prevTasks) => {
        return prevTasks.map((task) => {
          if (task.id === id) {
            const currentTime = Date.now();
            const elapsedTime = Math.floor((currentTime - startTime) / 1000);
            const newDuration = initialDuration + elapsedTime;

            if (newDuration >= task.total_time) {
              clearInterval(timer);
              return {
                ...task,
                duration: task.total_time,
                status: 'completed',
              };
            }

            if (newDuration < 0) {
              clearInterval(timer);
              return {
                ...task,
                duration: 0,
                status: 'pending',
              };
            }

            if (task.status === 'started') {
              return { ...task, duration: newDuration };
            }
          }
          return task;
        });
      });
    }, 1000);
    setTaskTimers((prevTimers) => ({ ...prevTimers, [id]: timer }));
  };

  const stopTaskTimer = (id: number) => {
    if (taskTimers[id]) {
      clearInterval(taskTimers[id]);
      setTaskTimers((timers) => {
        const newTimers = { ...timers };
        delete newTimers[id];
        return newTimers;
      });
    }
  };

  const fetchTasks = async () => {
    try {
      const fetchedTasks = await getTasks();
      if (Array.isArray(fetchedTasks)) {
        setTasks(fetchedTasks);
      } else {
        console.error('Fetched tasks are not in array format:', fetchedTasks);
        setTasks([]);
      }
    } catch (error) {
      console.error('Error fetching tasks:', error);
      setTasks([]);
    }
  };

  const handleCreateTask = async () => {
    let totalTimeInSeconds = newTotalTime;

    switch (timeUnit) {
      case 'minutes':
        totalTimeInSeconds = newTotalTime * 60;
        break;
      case 'hours':
        totalTimeInSeconds = newTotalTime * 3600;
        break;
      default:
        break;
    }

    try {
      const newTask = await createTask(newTaskName, totalTimeInSeconds);
      setTasks([...tasks, newTask]);
      setNewTaskName('');
    } catch (error) {
      console.error('Error creating task:', error);
    }
  };

  const handleUpdateTask = async (id: number, updatedTask: Partial<Task>) => {
    try {
      const task = await updateTask(id, updatedTask);
      setTasks(tasks.map((t) => (t.id === id ? task : t)));
    } catch (error) {
      console.error('Error updating task:', error);
    }
  };

  const handleDeleteTask = async (id: number) => {
    try {
      await deleteTask(id);
      setTasks(tasks.filter((t) => t.id !== id));
    } catch (error) {
      console.error('Error deleting task:', error);
    }
  };

  const handleStartTask = async (id: number) => {
    try {
      const task = await startTask(id);
      const startTime = Date.now();
      setTasks(
        tasks.map((t) =>
          t.id === id ? { ...task, start_time: new Date().toISOString() } : t,
        ),
      );
      startTaskTimer(id, 0, task.duration); // Start with 0 elapsed time but initial duration
    } catch (error) {
      console.error('Error starting task:', error);
    }
  };

  const handlePauseTask = async (id: number) => {
    try {
      const task = await pauseTask(id);
      stopTaskTimer(id); // Stop the timer first
      const currentDuration = tasks.find((t) => t.id === id)?.duration || 0;
      setTasks(
        tasks.map((t) =>
          t.id === id ? { ...task, duration: currentDuration } : t,
        ),
      );
    } catch (error) {
      console.error('Error pausing task:', error);
    }
  };

  const handleCompleteTask = async (id: number) => {
    try {
      const task = await completeTask(id);
      setTasks(tasks.map((t) => (t.id === id ? task : t)));
      stopTaskTimer(id);
    } catch (error) {
      console.error('Error completing task:', error);
    }
  };

  const handleResetTask = async (id: number) => {
    // 新增重置任务处理函数
    try {
      const task = await resetTask(id);
      setTasks(tasks.map((t) => (t.id === id ? task : t)));
      stopTaskTimer(id);
    } catch (error) {
      console.error('Error resetting task:', error);
    }
  };

  return (
    <div className={'container mx-auto p-4'}>
      <h1 className={'text-2xl font-bold mb-4'}>Welcome to TimerMe!</h1>
      <div className={'mb-4'}>
        <input
          type={'text'}
          value={newTaskName}
          onChange={(e) => setNewTaskName(e.target.value)}
          placeholder={'New Task Name'}
          className={'border p-2 mr-2'}
        />
        <input
          type={'number'}
          value={newTotalTime}
          onChange={(e) => setNewTotalTime(parseInt(e.target.value))}
          placeholder={'Duration(seconds)'}
          className={'border p-2 mr-2'}
        />

        <select
          value={timeUnit}
          onChange={(e) => setTimeUnit(e.target.value)}
          className='border p-2 mr-2'
        >
          <option value='seconds'>Seconds</option>
          <option value='minutes'>Minutes</option>
          <option value='hours'>Hours</option>
        </select>

        <button
          onClick={handleCreateTask}
          className={'bg-blue-500 text-white p-2 rounded'}
        >
          Add Task
        </button>
        <TaskList
          tasks={tasks}
          handleStartTask={handleStartTask}
          handlePauseTask={handlePauseTask}
          handleCompleteTask={handleCompleteTask}
          handleDeleteTask={handleDeleteTask}
          handleResetTask={handleResetTask} // 传递重置任务处理函数
        />
      </div>
    </div>
  );
}
