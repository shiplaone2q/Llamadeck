import { describe, expect, it } from 'vitest'
import { parseArgs } from '@/lib/parseArgs'

describe('parseArgs', () => {
  it('splits on whitespace', () => {
    expect(parseArgs('--delay=30s --port 8053')).toEqual(['--delay=30s', '--port', '8053'])
  })

  it('collapses repeated whitespace', () => {
    expect(parseArgs('  serve   --port   {port}  ')).toEqual(['serve', '--port', '{port}'])
  })
