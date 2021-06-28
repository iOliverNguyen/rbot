import IconMessenger from '../assets/IconMessenger.svelte';
import IconPage from '../assets/IconPage.svelte';
import IconPoint from '../assets/IconPoint.svelte';
import IconSend from '../assets/IconSend.svelte';
import IconTimer from '../assets/IconTimer.svelte';

import {newId, randN} from '../stores/util.js';

export type NodeType = string

export type NodeKind = 'trigger' | 'action' | 'branch' | 'invalid'

export type NodeTypeDecl = {
  icon: any
  title: string
  desc: string
  kind: NodeKind
}

export type Node = {
  id: number
  type: NodeType
  props: {}
  x: number
  y: number
}

export type Nodes = {
  nodes: Node[]
};

export const nodeTypes: Record<NodeType, NodeTypeDecl> = {
  webhook: {icon: IconPoint, kind: 'trigger', title: 'Webhook', desc: 'Trigger: by webhook callback'},
  message: {icon: IconMessenger, kind: 'trigger', title: 'Message', desc: 'Trigger: by response from Messenger'},
  timer: {icon: IconTimer, kind: 'trigger', title: 'Timer', desc: 'Trigger: after some delay'},
  send: {icon: IconSend, kind: 'action', title: 'Send Message', desc: 'Action: send message by Messenger'},
  sample: {icon: IconPage, kind: 'trigger', title: 'Sample', desc: 'An example of empty component'},
};

export const nodeTypeInvalid: NodeTypeDecl = {icon: undefined, kind: 'invalid', title: '', desc: ''};

export const nodeByType = (type: NodeType): NodeTypeDecl => nodeTypes[type];

export function newNodes(): Nodes {
  return {
    nodes: [
      newNode('webhook'),
      newNode('message'),
    ],
  };
}

export function newNode(type?: NodeType): Node {
  if (!type) {
    throw new Error(`invalid node type: ${type}`);
  }

  const node = {
    id: newId(),
    type: type,
    props: {},
    x: -200 + randN(400),
    y: -200 + randN(400),
  };
  return node;
}

export function insertNode(nodes: Nodes, node: Node) {
  nodes.nodes.push(node);
}
