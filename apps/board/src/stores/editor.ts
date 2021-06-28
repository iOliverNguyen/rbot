import {Writable, writable} from 'svelte/store';
import {PanelTabType} from '../core/editor.js';
import {newNodes, Nodes} from '../core/nodes.js';


export const nodes: Writable<Nodes> = writable(newNodes());

export const activePanelTab: Writable<PanelTabType> = writable('insert');
export const selectedNode: Writable<Node> = writable();
