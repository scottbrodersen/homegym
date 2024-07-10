export const select = {
  mounted: (el) => {
    const input =
      el.getElementsByTagName('input').length > 0
        ? el.getElementsByTagName('input')[0]
        : null;
    if (!!input && input.value == '0') {
      input.select();
    }
  },
};

export const focus = {
  mounted: (el) => {
    const input =
      el.getElementsByTagName('input').length > 0
        ? el.getElementsByTagName('input')[0]
        : null;

    if (
      (input &&
        input.hasAttribute('type') &&
        input.getAttribute('type') == 'number') ||
      (input && !isNaN(input.value))
    ) {
      if (input.value == '0') {
        input.focus();
      }
    } else {
      el.focus();
    }
  },
};
