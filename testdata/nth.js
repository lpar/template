function nth(o) {
  /* JS minification is pretty basic */
  return o + (['st','nd','rd'][(o+'').match(/1?\d\b/) - 1] || 'th');
}
