var random_store = {
    relations        : relations,
    suggested_dicts  : [],
    suggested_objects: [],
    keyword_dicts    : [],
    keyword_objects  : [],

    /**
     * Init
     */

    init() {
        this.suggested_dicts = this.random(9, random_dicts)
        this.suggested_objects = this.random(6, random_objects)
        this.keyword_dicts = this.random(5, random_dicts)
        this.keyword_objects = this.random(5, random_objects)
    },

    /**
     * random
     */

    random(amount, source) {
        var store = [];
        if (amount >= source.length) {
            return source;
        }
        while (store.length < amount) {
            var new_element = source[Math.floor(Math.random() * source.length)];
            while (store.find((v) => v.c !== undefined ? v.c === new_element.c : v.s === new_element.s)) {
                new_element = source[Math.floor(Math.random() * source.length)];
            }
            store = [...store, new_element];
        }
        return store
    }
}