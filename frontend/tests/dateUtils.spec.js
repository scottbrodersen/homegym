import { expect } from 'vitest';
import * as dateUtils from '../src/modules/dateUtils';

const testTimeStamp = 1715522280;
const testDate = new Date(testTimeStamp * 1000);

describe('dateUtils', () => {
  it('gets the date object', () => {
    const dateObj = dateUtils.dateFromSeconds(testTimeStamp);
    expect(dateObj).toEqual(testDate);
  });

  it('returns the current date for no timestamp', () => {
    const now = new Date();
    const dateObj = dateUtils.dateFromSeconds();
    expect(dateObj.getMilliseconds() >= now.getMilliseconds()).toBeTruthy;
  });

  it('gets the date string', () => {
    const dateStr = dateUtils.formatDate(testDate);
    expect(dateStr).toEqual('2024-05-12');
  });

  it('gets the time string', () => {
    const timeStr = dateUtils.formatTime(testDate);
    expect(timeStr).toEqual('9:58');
  });

  it('gets the date-time string', () => {
    const dateTimeStr = dateUtils.formatDateTime(testDate);
    expect(dateTimeStr).toEqual('2024-05-12 9:58');
  });

  it('gets the timestamp from date-time string', () => {
    const ts = dateUtils.stringToEpoch('2024-05-12 9:58');
    expect(ts).toEqual(testTimeStamp);
  });
});
