import type {Writable} from 'svelte/store';
import {writable} from 'svelte/store';

export interface XWritable<T> extends Writable<T> {
  // for getting back the value
  get(): T

  // for trigger rendering
  rfresh()
}

// this extends writable with get() and rfresh()
export function xwritable<T>(v: T): XWritable<T> {
  const w: any = writable(v);
  const _set = w.set;
  let _v = v;

  w.get = () => _v;
  w.set = (v) => {
    _v = v;
    _set.call(w, v);
  };
  w.rfresh = () => w.set(w.get());
  return w;
}

export const newId = (): number => new Date().getTime();

export const randN = (N: number) => Math.floor(Math.random() * N);
